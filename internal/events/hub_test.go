package events

import (
	"github.com/craftlane/craftlane/internal/protocol"
	"testing"
)

func TestSlowHubClientIsDisconnectedWithoutBlocking(t *testing.T) {
	hub := NewHub(1)
	client, cancel := hub.Subscribe()
	defer cancel()
	hub.Publish(protocol.Event{Sequence: 1})
	hub.Publish(protocol.Event{Sequence: 2})
	first, ok := <-client
	if !ok || first.Sequence != 1 {
		t.Fatalf("first=%#v ok=%v", first, ok)
	}
	if _, ok = <-client; ok {
		t.Fatal("slow client channel remained open")
	}
}
