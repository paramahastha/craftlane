package integration

import (
	"context"
	"github.com/craftlane/craftlane/internal/app"
	"github.com/craftlane/craftlane/internal/project"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/security"
	"github.com/craftlane/craftlane/tests/fixtures"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestActiveRuntimeEmitsCoalescedProjectChange(t *testing.T) {
	t.Setenv("CRAFTLANE_STATE_DIR", t.TempDir())
	dir := fixtures.GitRepo(t)
	commit := func(message string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "tracked"), []byte(message), 0600); err != nil {
			t.Fatal(err)
		}
		for _, args := range [][]string{{"config", "user.email", "test@example.com"}, {"config", "user.name", "Test"}, {"add", "."}, {"commit", "-m", message}} {
			cmd := exec.Command("git", args...)
			cmd.Dir = dir
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %v: %v (%s)", args, err, output)
			}
		}
	}
	commit("one")
	p, err := project.Detect(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := app.Start(context.Background(), p, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.Shutdown(context.Background())
	commit("two")
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		events, _ := runtime.Store.ReadEvents()
		for _, event := range events {
			if event.Type == "project.changed" {
				return
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("project.changed was not emitted")
}

func TestFoundationDetectsFixture(t *testing.T) {
	p, err := project.Detect(context.Background(), fixtures.GitRepo(t))
	if err != nil || p.ID == "" {
		t.Fatal(p, err)
	}
}

func TestFoundationRejectsNonGitAndRecoversStaleRecord(t *testing.T) {
	if _, err := project.Detect(context.Background(), t.TempDir()); err == nil {
		t.Fatal("expected non-Git directory rejection")
	}
	t.Setenv("CRAFTLANE_STATE_DIR", t.TempDir())
	p, err := project.Detect(context.Background(), fixtures.GitRepo(t))
	if err != nil {
		t.Fatal(err)
	}
	store, err := app.StorageForProject(p)
	if err != nil {
		t.Fatal(err)
	}
	if err = store.SaveRecord(protocol.RuntimeRecord{RuntimeID: "stale", ProjectID: p.ID, ProjectRoot: p.Root, Origin: "http://127.0.0.1:1", State: "ready"}); err != nil {
		t.Fatal(err)
	}
	if err = store.SaveControl("obsolete"); err != nil {
		t.Fatal(err)
	}
	record, live, err := app.RuntimeStatus(context.Background(), p)
	if err != nil || live || record.State != "stale" {
		t.Fatalf("record=%#v live=%v err=%v", record, live, err)
	}
	if _, err = store.LoadControl(); !os.IsNotExist(err) {
		t.Fatalf("obsolete control was not removed: %v", err)
	}
}

func TestFoundationLoopbackOnlyAndAtomicFiles(t *testing.T) {
	runtime, _ := startRuntime(t)
	host, _, err := net.SplitHostPort(runtime.Listener.Addr().String())
	if err != nil || host != "127.0.0.1" {
		t.Fatalf("listener=%s err=%v", runtime.Listener.Addr(), err)
	}
	if _, err := os.Stat(filepath.Join(runtime.Store.Dir, "runtime.json.tmp")); !os.IsNotExist(err) {
		t.Fatalf("temporary record remains: %v", err)
	}
	req, _ := http.NewRequest(http.MethodGet, runtime.Record.Origin+"/api/v1/health", nil)
	req.Host = "evil.test"
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid Host status=%d", resp.StatusCode)
	}
	outside := t.TempDir()
	if security.ApprovedPath(runtime.Project.Root, outside) {
		t.Fatal("out-of-project path approved")
	}
}
