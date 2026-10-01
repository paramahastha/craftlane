package events

import (
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/storage"
	"os"
	"strings"
	"testing"
)

func TestReducerPersistsBeforeBroadcastAndRecoversSequence(t *testing.T) {
	store := storage.New(t.TempDir())
	reducer := NewReducer(store, protocol.Project{ID: "project"}, protocol.Runtime{ID: "runtime", State: "ready"})
	live, cancel, _ := reducer.Subscribe()
	defer cancel()
	event, err := reducer.Append("project.changed", map[string]any{"token": "secret-value", "nested": map[string]string{"password": "also-secret"}})
	if err != nil {
		t.Fatal(err)
	}
	if got := <-live; got.Sequence != event.Sequence {
		t.Fatalf("broadcast sequence %d, want %d", got.Sequence, event.Sequence)
	}
	data, err := os.ReadFile(store.Dir + "/events.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "secret-value") || strings.Contains(string(data), "also-secret") {
		t.Fatal("raw secret reached event storage")
	}
	recovered := NewReducer(store, protocol.Project{ID: "project"}, protocol.Runtime{ID: "runtime"})
	next, err := recovered.Append("runtime.ready", map[string]any{})
	if err != nil || next.Sequence != event.Sequence+1 {
		t.Fatalf("recovered append = %#v, %v", next, err)
	}
	if snapshot, err := store.LoadSnapshot(); err != nil || snapshot.Sequence != next.Sequence {
		t.Fatalf("snapshot = %#v, %v", snapshot, err)
	}
}

func TestRetentionHasMinimumAndRequestsResync(t *testing.T) {
	store := storage.New(t.TempDir())
	reducer := NewReducerWithRetention(store, protocol.Project{ID: "p"}, protocol.Runtime{}, 2)
	if reducer.Retention != storage.MinimumRetention {
		t.Fatalf("retention = %d", reducer.Retention)
	}
	for i := 0; i < storage.MinimumRetention+2; i++ {
		if _, err := reducer.Append("project.changed", map[string]any{}); err != nil {
			t.Fatal(err)
		}
	}
	status, replay, _ := reducer.Replay(0, reducer.Sequence)
	if status != "resync_required" || len(replay) != 0 {
		t.Fatalf("status=%s replay=%d", status, len(replay))
	}
}
