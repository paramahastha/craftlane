package codex

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/craftlane/craftlane/internal/protocol"
)

func TestMalformedConfigIsNeverOverwritten(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".codex", "hooks.json")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("not-json"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := New(root).Apply(); err == nil {
		t.Fatal("malformed config was accepted")
	}
	b, _ := os.ReadFile(path)
	if string(b) != "not-json" {
		t.Fatal("malformed config changed")
	}
}

func TestNormalizeDoesNotTreatTurnStopAsSessionTerminal(t *testing.T) {
	e, err := New(t.TempDir()).Normalize(protocol.ObservationEnvelope{Version: 1, Provider: "codex", ProjectID: "p", Kind: "turn.completed", SessionKey: "s"})
	if err != nil || e.Kind != "turn.updated" {
		t.Fatalf("normalized = %#v, %v", e, err)
	}
}
