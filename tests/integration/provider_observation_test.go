package integration

import (
	"testing"

	"github.com/craftlane/craftlane/internal/events"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/storage"
)

func TestProviderObservationKeepsParentRelationSafe(t *testing.T) {
	p := protocol.Project{ID: "p"}
	r := events.NewReducer(storage.New(t.TempDir()), p, protocol.Runtime{})
	agent, changed, err := r.ApplyObservation(protocol.ObservationEnvelope{Version: 1, Provider: "opencode", ProjectID: "p", Kind: "session.busy", SessionKey: "s", AgentKey: "a", ParentAgentKey: "parent"}, protocol.EventSource{Name: "opencode", AdapterVersion: "v1"})
	if err != nil || !changed || agent.ParentAgentID == "" {
		t.Fatalf("agent = %#v, %v, %v", agent, changed, err)
	}
	if len(r.Events) != 1 || r.Events[0].Source.Name != "opencode" {
		t.Fatalf("events = %#v", r.Events)
	}
}
