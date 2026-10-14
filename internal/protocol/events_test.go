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

func TestSnapshotAgentAndProviderFieldsAreAdditive(t *testing.T) {
	data := []byte(`{"sequence":2,"project":{"id":"p","name":"p","root":"{{PROJECT_ROOT}}"},"runtime":{"id":"r","origin":"http://127.0.0.1:1","protocolVersion":1,"startedAt":"2026-10-02T00:00:00Z","state":"ready"},"diagnostics":[],"agents":[{"id":"a","provider":"claude","projectId":"p","displayName":"agent","state":"working","stateSource":"reported","lastObservedAt":"2026-10-02T00:00:00Z"}],"providerStatuses":[{"provider":"codex","setupState":"unavailable"}]}`)
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Agents) != 1 || snapshot.Agents[0].Provider != ProviderClaude || len(snapshot.ProviderStatuses) != 1 {
		t.Fatalf("snapshot = %#v", snapshot)
	}
}
