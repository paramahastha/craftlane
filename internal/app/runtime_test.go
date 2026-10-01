package app

import (
	"context"
	"github.com/craftlane/craftlane/internal/project"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/storage"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRuntimeServesHealth(t *testing.T) {
	p := project.New(t.TempDir(), t.TempDir())
	r, err := Start(context.Background(), p, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Shutdown(context.Background())
	req := httptest.NewRequest(http.MethodGet, r.Record.Origin+"/api/v1/health", nil)
	req.Host = r.Record.Origin[len("http://"):]
	w := httptest.NewRecorder()
	r.HTTP.Handler.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("got %d", w.Code)
	}
}

func TestStartupFailureReleasesListenerAndInvalidatesState(t *testing.T) {
	stateFile := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(stateFile, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CRAFTLANE_STATE_DIR", stateFile)
	probe, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := probe.Addr().(*net.TCPAddr).Port
	_ = probe.Close()
	if _, err = Start(context.Background(), project.New(t.TempDir(), t.TempDir()), port); err == nil {
		t.Fatal("expected startup failure")
	}

	t.Setenv("CRAFTLANE_STATE_DIR", t.TempDir())
	runtime, err := Start(context.Background(), project.New(t.TempDir(), t.TempDir()), port)
	if err != nil {
		t.Fatalf("listener was not released after failed startup: %v", err)
	}
	_ = runtime.Shutdown(context.Background())
}

func TestDetachedStartupCleanupStopsChildAndInvalidatesControl(t *testing.T) {
	store := storage.New(t.TempDir())
	if err := store.SaveRecord(protocol.RuntimeRecord{RuntimeID: "runtime", State: "starting"}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveControl("secret"); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestDetachedCleanupChild")
	cmd.Env = append(os.Environ(), "CRAFTLANE_TEST_CHILD=1")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	cleanupDetachedStart(cmd, store, "runtime")

	if _, err := store.LoadControl(); !os.IsNotExist(err) {
		t.Fatalf("control secret remains after failed detached startup: %v", err)
	}
	record, err := store.LoadRecord()
	if err != nil {
		t.Fatal(err)
	}
	if record.State != "stopped" {
		t.Fatalf("record state = %q, want stopped", record.State)
	}
}

func TestDetachedCleanupChild(t *testing.T) {
	if os.Getenv("CRAFTLANE_TEST_CHILD") != "1" {
		return
	}
	select {}
}
