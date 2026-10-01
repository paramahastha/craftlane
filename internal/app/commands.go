package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/craftlane/craftlane/internal/diagnostics"
	"github.com/craftlane/craftlane/internal/project"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/server"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func VerifyRuntime(ctx context.Context, p protocol.Project, record protocol.RuntimeRecord) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, record.Origin+"/api/v1/health", nil)
	if err != nil {
		return err
	}
	request.Host = request.URL.Host
	response, err := (&http.Client{Timeout: time.Second}).Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("health returned %s", response.Status)
	}
	var health protocol.Health
	if err = json.NewDecoder(response.Body).Decode(&health); err != nil {
		return err
	}
	if health.Runtime.ID != record.RuntimeID || health.Project.ID != p.ID || health.Runtime.Origin != record.Origin {
		return fmt.Errorf("runtime health identity does not match the stored record")
	}
	return nil
}

func RuntimeStatus(ctx context.Context, p protocol.Project) (protocol.RuntimeRecord, bool, error) {
	store, err := StorageForProject(p)
	if err != nil {
		return protocol.RuntimeRecord{}, false, err
	}
	record, err := store.LoadRecord()
	if os.IsNotExist(err) {
		return protocol.RuntimeRecord{ProjectID: p.ID, ProjectRoot: p.Root, State: "stopped"}, false, nil
	}
	if err != nil {
		return protocol.RuntimeRecord{}, false, fmt.Errorf("read runtime state: %w", err)
	}
	if VerifyRuntime(ctx, p, record) == nil {
		return record, true, nil
	}
	if record.State != "stopped" {
		record.State = "stale"
		_ = store.SaveRecord(record)
	}
	_ = store.ClearControl()
	return record, false, nil
}

func StopProject(ctx context.Context, p protocol.Project) error {
	store, err := StorageForProject(p)
	if err != nil {
		return err
	}
	record, live, err := RuntimeStatus(ctx, p)
	if err != nil {
		return err
	}
	if !live {
		return fmt.Errorf("no matching live runtime (state: %s); run craftlane start", record.State)
	}
	token, err := store.LoadControl()
	if err != nil {
		return fmt.Errorf("runtime control is unavailable; restart the runtime")
	}
	u, err := url.Parse(record.Origin + "/api/v1/runtime/stop")
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return err
	}
	request.Host = u.Host
	request.Header.Set("X-Craftlane-Token", token)
	response, err := (&http.Client{Timeout: 3 * time.Second}).Do(request)
	if err != nil {
		return fmt.Errorf("runtime is not reachable: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("runtime rejected stop request (%s)", response.Status)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if VerifyRuntime(ctx, p, record) != nil {
			return nil
		}
		time.Sleep(25 * time.Millisecond)
	}
	return fmt.Errorf("runtime did not stop within 5 seconds")
}

func OpenProject(ctx context.Context, p protocol.Project) error {
	record, live, err := RuntimeStatus(ctx, p)
	if err != nil {
		return err
	}
	if !live {
		return fmt.Errorf("no live runtime; run craftlane start")
	}
	return OpenBrowser(record.Origin)
}

func Doctor(ctx context.Context, p protocol.Project) ([]protocol.DiagnosticCheck, bool) {
	_, live, _ := RuntimeStatus(ctx, p)
	checks := diagnostics.Detect(p.Root, live, server.AssetReady())
	ok := true
	for _, check := range checks {
		if check.Status == "fail" {
			ok = false
		}
	}
	return checks, ok
}

func Serve(root string, port int, runtimeID string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	p, err := project.Detect(ctx, root)
	if err != nil {
		return err
	}
	runtime, err := StartWithID(ctx, p, port, runtimeID)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		shutdownCtx, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		return runtime.Shutdown(shutdownCtx)
	case <-runtime.Done():
		return nil
	}
}
