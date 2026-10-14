package opencode

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/craftlane/craftlane/internal/protocol"
)

func TestOnlyCraftlaneOwnedPluginCanBeRemoved(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".opencode", "plugins", "craftlane-observability.ts")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("export const userPlugin = true;\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := New(root).Remove(); err == nil {
		t.Fatal("unowned plugin was removable")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("unowned plugin disappeared")
	}
	if result, err := New(root).Apply(); err == nil || result.Changed {
		t.Fatal("unowned plugin was overwritten")
	}
}

func TestNormalizePreservesSuppliedParentLineage(t *testing.T) {
	e, err := New(t.TempDir()).Normalize(protocol.ObservationEnvelope{Version: 1, Provider: "opencode", ProjectID: "p", Kind: "session.busy", SessionKey: "s", AgentKey: "a", ParentAgentKey: "parent"})
	if err != nil || e.ParentAgentKey != "parent" {
		t.Fatalf("normalized = %#v, %v", e, err)
	}
}
