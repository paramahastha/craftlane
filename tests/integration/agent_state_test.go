package integration

import (
	"testing"
	"time"

	"github.com/craftlane/craftlane/internal/events"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/storage"
)

func TestInferredIdleIsDistinctFromReportedPermission(t *testing.T) {
	r := events.NewReducer(storage.New(t.TempDir()), protocol.Project{ID: "p"}, protocol.Runtime{})
	_, _, err := r.ApplyObservation(protocol.ObservationEnvelope{Version: 1, Provider: "claude", ProjectID: "p", Kind: "session.started", SessionKey: "s"}, protocol.EventSource{Name: "claude"})
	if err != nil {
		t.Fatal(err)
	}
	changed, err := r.InferIdle(time.Now().Add(time.Minute), time.Second)
	if err != nil || len(changed) != 1 || changed[0].StateSource != "inferred" {
		t.Fatalf("idle = %#v, %v", changed, err)
	}
	permission, _, err := r.ApplyObservation(protocol.ObservationEnvelope{Version: 1, Provider: "claude", ProjectID: "p", Kind: "permission.requested", SessionKey: "s"}, protocol.EventSource{Name: "claude"})
	if err != nil || permission.State != protocol.AgentWaitingPermission || permission.StateSource != "reported" {
		t.Fatalf("permission = %#v, %v", permission, err)
	}
}
