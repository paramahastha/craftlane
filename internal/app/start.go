package app

import (
	"context"
	"fmt"
	"github.com/craftlane/craftlane/internal/project"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/storage"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const defaultPort = 7331

func StartDetached(ctx context.Context, dir string, port int, explicitPort, openBrowser bool) (protocol.RuntimeRecord, error) {
	p, err := project.Detect(ctx, dir)
	if err != nil {
		return protocol.RuntimeRecord{}, err
	}
	store, err := StorageForProject(p)
	if err != nil {
		return protocol.RuntimeRecord{}, err
	}
	if record, err := store.LoadRecord(); err == nil {
		if VerifyRuntime(ctx, p, record) == nil {
			if openBrowser {
				_ = OpenBrowser(record.Origin)
			}
			return record, nil
		}
		record.State = "stale"
		_ = store.SaveRecord(record)
		_ = store.ClearControl()
	}
	if port == 0 {
		port = defaultPort
	}
	selected, err := selectPort(port, explicitPort)
	if err != nil {
		return protocol.RuntimeRecord{}, err
	}
	runtimeID := storage.RandomSecret()[:16]
	executable, err := os.Executable()
	if err != nil {
		return protocol.RuntimeRecord{}, err
	}
	cmd := exec.Command(executable, "serve", "--root", p.Root, "--port", strconv.Itoa(selected), "--runtime-id", runtimeID)
	cmd.Env = os.Environ()
	detachProcess(cmd)
	if err = cmd.Start(); err != nil {
		return protocol.RuntimeRecord{}, fmt.Errorf("start runtime process: %w", err)
	}
	cleanup := func() { cleanupDetachedStart(cmd, store, runtimeID) }
	deadline := time.Now().Add(10 * time.Second)
	var record protocol.RuntimeRecord
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			cleanup()
			return protocol.RuntimeRecord{}, ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
		record, err = store.LoadRecord()
		if err == nil && record.RuntimeID == runtimeID && VerifyRuntime(ctx, p, record) == nil {
			_ = cmd.Process.Release()
			if openBrowser {
				_ = OpenBrowser(record.Origin)
			}
			return record, nil
		}
	}
	cleanup()
	return protocol.RuntimeRecord{}, fmt.Errorf("runtime did not become ready on 127.0.0.1:%d within 10 seconds; check that the port is available and retry", selected)
}

func cleanupDetachedStart(cmd *exec.Cmd, store storage.Store, runtimeID string) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}
	if current, loadErr := store.LoadRecord(); loadErr == nil && current.RuntimeID == runtimeID {
		current.State = "stopped"
		_ = store.SaveRecord(current)
	}
	_ = store.ClearControl()
}

func selectPort(port int, explicit bool) (int, error) {
	last := port
	if !explicit {
		last = port + 20
	}
	for candidate := port; candidate <= last; candidate++ {
		listener, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", candidate))
		if err == nil {
			_ = listener.Close()
			return candidate, nil
		}
	}
	if explicit {
		return 0, fmt.Errorf("requested port %d is unavailable; choose another port with --port", port)
	}
	return 0, fmt.Errorf("ports %d-%d are unavailable; free a port or use --port", port, last)
}
