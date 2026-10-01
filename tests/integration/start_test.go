package integration

import (
	"github.com/craftlane/craftlane/tests/fixtures"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestStartModesConflictAndTwentyCleanStarts(t *testing.T) {
	bin := buildBinary(t)
	projectDir := fixtures.GitRepo(t)
	state := t.TempDir()
	t.Setenv("BROWSER", bin)
	if output, err := runCLI(t, bin, projectDir, state); err != nil || !strings.Contains(output, "http://127.0.0.1:") {
		t.Fatalf("default start: output=%q err=%v", output, err)
	}
	if _, err := runCLI(t, bin, projectDir, state, "stop"); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		started := time.Now()
		output, err := runCLI(t, bin, projectDir, state, "start", "--no-open")
		if err != nil || !strings.Contains(output, "http://127.0.0.1:") || time.Since(started) > 10*time.Second {
			t.Fatalf("start %d: output=%q err=%v duration=%s", i, output, err, time.Since(started))
		}
		if _, err = runCLI(t, bin, projectDir, state, "stop"); err != nil {
			t.Fatalf("stop %d: %v", i, err)
		}
	}
	if output, err := runCLI(t, bin, projectDir, state, "start", "--headless"); err != nil || !strings.Contains(output, "127.0.0.1") {
		t.Fatalf("headless: %q %v", output, err)
	}
	if _, err := runCLI(t, bin, projectDir, state, "stop"); err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	defer listener.Close()
	output, err := runCLI(t, bin, projectDir, state, "start", "--no-open", "--port", strconv.Itoa(port))
	if err == nil || !strings.Contains(output, "choose another port") {
		t.Fatalf("conflict output=%q err=%v", output, err)
	}
}
