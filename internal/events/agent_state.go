package events

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/craftlane/craftlane/internal/protocol"
)

func ObservationFingerprint(e protocol.ObservationEnvelope) string {
	key := strings.Join([]string{e.Provider, e.ProjectID, e.Kind, e.SessionKey, e.AgentKey, e.ParentAgentKey, e.ToolKey, e.ToolName, e.Outcome, e.ReportedAt.UTC().Format(time.RFC3339Nano)}, "\x00")
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func ObservationAgentID(e protocol.ObservationEnvelope) string {
	key := e.AgentKey
	if key == "" {
		key = e.SessionKey
	}
	h := sha256.Sum256([]byte(e.Provider + "\x00" + e.ProjectID + "\x00" + key))
	return hex.EncodeToString(h[:])
}

func agentDisplayName(e protocol.ObservationEnvelope) string {
	if e.AgentKey != "" {
		return e.Provider + " agent " + e.AgentKey
	}
	return e.Provider + " session " + e.SessionKey
}

func stateForObservation(kind string) (protocol.AgentState, string, bool) {
	switch {
	case strings.HasPrefix(kind, "tool."):
		return protocol.AgentWorking, "reported", true
	case strings.Contains(kind, "permission"):
		return protocol.AgentWaitingPermission, "reported", true
	case strings.HasSuffix(kind, ".idle"):
		return protocol.AgentIdle, "reported", true
	case strings.HasSuffix(kind, ".waiting"):
		return protocol.AgentWaiting, "reported", true
	case strings.Contains(kind, "failed"), strings.HasSuffix(kind, ".error"):
		return protocol.AgentFailed, "reported", true
	case strings.Contains(kind, "completed"), strings.Contains(kind, "stopped"), strings.HasSuffix(kind, ".ended"), strings.HasSuffix(kind, ".end"):
		return protocol.AgentCompleted, "reported", true
	case strings.Contains(kind, "started"), strings.Contains(kind, "updated"), strings.Contains(kind, "activity"), strings.Contains(kind, "busy"):
		return protocol.AgentWorking, "reported", true
	default:
		return protocol.AgentUnknown, "reported", false
	}
}

func applyObservation(agent protocol.Agent, e protocol.ObservationEnvelope) (protocol.Agent, bool, error) {
	state, source, supported := stateForObservation(e.Kind)
	if !supported {
		return agent, false, fmt.Errorf("unsupported observation kind %q", e.Kind)
	}
	if agent.ID == "" {
		agent = protocol.Agent{ID: ObservationAgentID(e), Provider: e.Provider, ProjectID: e.ProjectID, SessionID: e.SessionKey, DisplayName: agentDisplayName(e), State: protocol.AgentUnknown, StateSource: source, PriorHistoryUnavailable: !strings.HasSuffix(e.Kind, "started")}
	}
	if agent.State == protocol.AgentCompleted || agent.State == protocol.AgentFailed {
		return agent, false, nil
	}
	if e.ParentAgentKey != "" {
		parent := e
		parent.AgentKey = e.ParentAgentKey
		agent.ParentAgentID = ObservationAgentID(parent)
	}
	if e.ToolName != "" || e.ToolKey != "" {
		toolState := "started"
		if strings.Contains(e.Kind, "completed") {
			toolState = "completed"
		}
		if strings.Contains(e.Kind, "failed") {
			toolState = "failed"
		}
		agent.CurrentTool = &protocol.ToolActivity{ID: e.ToolKey, Name: e.ToolName, State: toolState, Outcome: e.Outcome}
	}
	agent.State, agent.StateSource = state, source
	agent.StaleAfterRestart = false
	agent.LastObservedAt = time.Now().UTC()
	if !e.ReportedAt.IsZero() {
		agent.ReportedAt = e.ReportedAt
	}
	return agent, true, nil
}
