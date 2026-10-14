package integration

import (
	"testing"

	"github.com/craftlane/craftlane/internal/app"
	"github.com/craftlane/craftlane/internal/protocol"
)

func TestObserveSetupPreviewAndRemovalAreProviderScoped(t *testing.T) {
	p := protocol.Project{ID: "p", Root: t.TempDir()}
	preview, err := app.ObserveSetup(p, "claude", false)
	if err != nil || len(preview.Preview.Changes) != 1 {
		t.Fatalf("preview = %#v, %v", preview, err)
	}
	result, err := app.ObserveSetup(p, "claude", true)
	if err != nil || !result.Changed {
		t.Fatalf("apply = %#v, %v", result, err)
	}
	removed, err := app.ObserveRemove(p, "claude", true)
	if err != nil || !removed.Changed {
		t.Fatalf("remove = %#v, %v", removed, err)
	}
	if _, err := app.ObserveSetup(p, "unknown", false); err == nil {
		t.Fatal("unknown provider accepted")
	}
}
