package integration

import (
	"github.com/craftlane/craftlane/tests/fixtures"
	"strings"
	"testing"
)

func TestBuiltBinaryRunsWithoutWebTooling(t *testing.T) {
	t.Setenv("HTTP_PROXY", "http://127.0.0.1:1")
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:1")
	bin := buildBinary(t)
	projectDir := fixtures.GitRepo(t)
	state := t.TempDir()
	output, err := runCLI(t, bin, projectDir, state, "start", "--no-open")
	if err != nil || !strings.Contains(output, "http://127.0.0.1:") {
		t.Fatalf("start=%q err=%v", output, err)
	}
	if _, err = runCLI(t, bin, projectDir, state, "stop"); err != nil {
		t.Fatal(err)
	}
}
