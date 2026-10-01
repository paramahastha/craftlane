package protocol

import (
	"encoding/json"
	"testing"
)

func TestEventEnvelopeIgnoresUnknownOptionalFields(t *testing.T) {
	data := []byte(`{"schemaVersion":1,"eventId":"e","sequence":1,"timestamp":"2026-10-01T00:00:00Z","projectId":"p","source":{"name":"core","future":"ok"},"type":"project.changed","payload":{},"future":"ok"}`)
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatal(err)
	}
	if event.SchemaVersion != 1 || event.EventID != "e" || event.Source.Name != "core" {
		t.Fatalf("unexpected event: %#v", event)
	}
}
