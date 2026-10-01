package integration

import (
	"github.com/craftlane/craftlane/tests/fixtures"
	"strings"
	"testing"
)

func TestCLILifecycleStatusRestartStop(t *testing.T) {
	bin := buildBinary(t)
	t.Setenv("BROWSER", bin)
	projectDir := fixtures.GitRepo(t)
	state := t.TempDir()
	if output, err := runCLI(t, bin, projectDir, state, "status"); err != nil || !strings.Contains(output, "State: stopped") {
		t.Fatalf("initial status: %q %v", output, err)
	}
	if _, err := runCLI(t, bin, projectDir, state, "start", "--no-open"); err != nil {
		t.Fatal(err)
	}
	if output, err := runCLI(t, bin, projectDir, state, "status"); err != nil || !strings.Contains(output, "State: ready") {
		t.Fatalf("running status: %q %v", output, err)
	}
	if _, err := runCLI(t, bin, projectDir, state, "open"); err != nil {
		t.Fatalf("open: %v", err)
	}
	if output, err := runCLI(t, bin, projectDir, state, "restart"); err != nil || !strings.Contains(output, "restarted") {
		t.Fatalf("restart: %q %v", output, err)
	}
	if _, err := runCLI(t, bin, projectDir, state, "stop"); err != nil {
		t.Fatal(err)
	}
	if output, err := runCLI(t, bin, projectDir, state, "status"); err != nil || !strings.Contains(output, "State: stopped") {
		t.Fatalf("stopped status: %q %v", output, err)
	}
}
