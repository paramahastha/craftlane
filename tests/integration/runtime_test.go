package integration

import (
	"context"
	"encoding/json"
	"github.com/craftlane/craftlane/internal/app"
	"github.com/craftlane/craftlane/internal/project"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/tests/fixtures"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func startRuntime(t *testing.T) (*app.Runtime, protocol.Project) {
	t.Helper()
	t.Setenv("CRAFTLANE_STATE_DIR", t.TempDir())
	p, err := project.Detect(context.Background(), fixtures.GitRepo(t))
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := app.Start(context.Background(), p, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = runtime.Shutdown(context.Background()) })
	return runtime, p
}

func getJSON(t *testing.T, origin, path string, target any) int {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, origin+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = req.URL.Host
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if target != nil && resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			t.Fatal(err)
		}
	}
	return resp.StatusCode
}

func buildBinary(t *testing.T) string {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	name := "craftlane"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(t.TempDir(), name)
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/craftlane")
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build binary: %v\n%s", err, output)
	}
	return bin
}

func runCLI(t *testing.T, bin, dir, state string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "CRAFTLANE_STATE_DIR="+state)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}
