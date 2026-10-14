package events

import (
	"testing"

	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/storage"
)

func observation(project, kind string) protocol.ObservationEnvelope {
	return protocol.ObservationEnvelope{Version: 1, Provider: protocol.ProviderClaude, ProjectID: project, Kind: kind, SessionKey: "session-1", AgentKey: "agent-1"}
}

func TestReducerObservationStateGuardsAndDuplicateSuppression(t *testing.T) {
	projectID := "project"
	r := NewReducer(storage.New(t.TempDir()), protocol.Project{ID: projectID}, protocol.Runtime{})
	agent, changed, err := r.ApplyObservation(observation(projectID, "session.started"), protocol.EventSource{Name: "claude", AdapterVersion: "v1"})
	if err != nil || !changed || agent.State != protocol.AgentWorking {
		t.Fatalf("start = %#v, %v, %v", agent, changed, err)
	}
	if _, changed, err = r.ApplyObservation(observation(projectID, "session.started"), protocol.EventSource{Name: "claude"}); err != nil || changed {
		t.Fatalf("duplicate = %v, %v", changed, err)
	}
	permission := observation(projectID, "permission.requested")
	permission.ToolName = "read_file"
	agent, changed, err = r.ApplyObservation(permission, protocol.EventSource{Name: "claude"})
	if err != nil || !changed || agent.State != protocol.AgentWaitingPermission || agent.StateSource != "reported" {
		t.Fatalf("permission = %#v, %v, %v", agent, changed, err)
	}
	completed := observation(projectID, "session.completed")
	agent, _, err = r.ApplyObservation(completed, protocol.EventSource{Name: "claude"})
	if err != nil || agent.State != protocol.AgentCompleted {
		t.Fatalf("complete = %#v, %v", agent, err)
	}
	later := observation(projectID, "tool.started")
	agent, changed, err = r.ApplyObservation(later, protocol.EventSource{Name: "claude"})
	if err != nil || changed || agent.State != protocol.AgentCompleted {
		t.Fatalf("terminal guard = %#v, %v, %v", agent, changed, err)
	}
}

func TestReducerObservationMarksLateHistoryUnavailable(t *testing.T) {
	r := NewReducer(storage.New(t.TempDir()), protocol.Project{ID: "p"}, protocol.Runtime{})
	agent, changed, err := r.ApplyObservation(observation("p", "tool.completed"), protocol.EventSource{Name: "claude"})
	if err != nil || !changed || !agent.PriorHistoryUnavailable {
		t.Fatalf("late observation = %#v, %v, %v", agent, changed, err)
	}
}
