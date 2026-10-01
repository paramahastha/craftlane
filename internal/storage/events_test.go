package storage

import (
	"github.com/craftlane/craftlane/internal/protocol"
	"os"
	"testing"
)

func TestJSONLEvents(t *testing.T) {
	s := New(t.TempDir())
	if err := s.AppendEvent(protocol.Event{Sequence: 1, EventID: "e"}); err != nil {
		t.Fatal(err)
	}
	v, err := s.ReadEvents()
	if err != nil || len(v) != 1 || v[0].EventID != "e" {
		t.Fatalf("%v %v", v, err)
	}
}

func TestNormalizeRetention(t *testing.T) {
	if got := NormalizeRetention(0); got != DefaultRetention {
		t.Fatalf("default = %d", got)
	}
	if got := NormalizeRetention(1); got != MinimumRetention {
		t.Fatalf("minimum = %d", got)
	}
}

func TestEventsUseBoundedAppendOnlySegments(t *testing.T) {
	s := New(t.TempDir())
	for i := 1; i <= segmentSize+1; i++ {
		if err := s.AppendEvent(protocol.Event{Sequence: uint64(i), EventID: "e"}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(s.path("events-000002.jsonl")); err != nil {
		t.Fatalf("segment rollover: %v", err)
	}
	events, err := s.ReadEvents()
	if err != nil || len(events) != segmentSize+1 {
		t.Fatalf("read segmented events = %d, %v", len(events), err)
	}
}

func TestRemoveEventRollsBackLatestAppend(t *testing.T) {
	s := New(t.TempDir())
	if err := s.AppendEvent(protocol.Event{Sequence: 1, EventID: "e"}); err != nil {
		t.Fatal(err)
	}
	if err := s.RemoveEvent(1); err != nil {
		t.Fatal(err)
	}
	events, err := s.ReadEvents()
	if err != nil || len(events) != 0 {
		t.Fatalf("rolled back events = %d, %v", len(events), err)
	}
}

func TestTruncateEventsAfterSnapshotCursor(t *testing.T) {
	s := New(t.TempDir())
	for i := 1; i <= 2; i++ {
		if err := s.AppendEvent(protocol.Event{Sequence: uint64(i), EventID: "e"}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.TruncateEventsAfter(1); err != nil {
		t.Fatal(err)
	}
	events, err := s.ReadEvents()
	if err != nil || len(events) != 1 || events[0].Sequence != 1 {
		t.Fatalf("recovered events = %#v, %v", events, err)
	}
}
