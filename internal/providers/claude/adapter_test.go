package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/craftlane/craftlane/internal/protocol"
)

func TestSetupMergesAndRemovesOnlyCraftlaneHook(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".claude", "settings.local.json")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"hooks":{"SessionStart":[{"command":"user-hook"}]},"theme":"dark"}`), 0600); err != nil {
		t.Fatal(err)
	}
	a := New(root)
	preview, err := a.Preview()
	if err != nil || preview.AlreadyApplied {
		t.Fatalf("preview = %#v, %v", preview, err)
	}
	if result, err := a.Apply(); err != nil || !result.Changed {
		t.Fatalf("apply = %#v, %v", result, err)
	}
	if result, err := a.Apply(); err != nil || result.Changed {
		t.Fatalf("idempotent apply = %#v, %v", result, err)
	}
	if result, err := a.Remove(); err != nil || !result.Changed {
		t.Fatalf("remove = %#v, %v", result, err)
	}
	var object map[string]any
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &object); err != nil || object["theme"] != "dark" {
		t.Fatalf("preserved config = %#v, %v", object, err)
	}
}

func TestNormalizeKeepsOnlySafeLifecycleEnvelope(t *testing.T) {
	e, err := New(t.TempDir()).Normalize(protocol.ObservationEnvelope{Version: 1, Provider: "claude", ProjectID: "p", Kind: "tool.started", SessionKey: "s", AgentKey: "a", ToolName: "read_file"})
	if err != nil || e.ToolName != "read_file" {
		t.Fatalf("normalized = %#v, %v", e, err)
	}
}
