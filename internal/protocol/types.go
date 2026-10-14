package protocol

import "time"

type Capabilities map[string]string

const (
	ProviderClaude   = "claude"
	ProviderCodex    = "codex"
	ProviderOpenCode = "opencode"
)

type AgentState string

const (
	AgentWorking           AgentState = "working"
	AgentWaiting           AgentState = "waiting"
	AgentWaitingPermission AgentState = "waiting_permission"
	AgentIdle              AgentState = "idle"
	AgentCompleted         AgentState = "completed"
	AgentFailed            AgentState = "failed"
	AgentUnknown           AgentState = "unknown"
)

type ToolActivity struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Outcome string `json:"outcome,omitempty"`
}

type Agent struct {
	ID                      string        `json:"id"`
	Provider                string        `json:"provider"`
	ProjectID               string        `json:"projectId"`
	SessionID               string        `json:"sessionId,omitempty"`
	DisplayName             string        `json:"displayName"`
	ParentAgentID           string        `json:"parentAgentId,omitempty"`
	State                   AgentState    `json:"state"`
	StateSource             string        `json:"stateSource"`
	CurrentTool             *ToolActivity `json:"currentTool,omitempty"`
	Task                    string        `json:"task,omitempty"`
	LastObservedAt          time.Time     `json:"lastObservedAt"`
	ReportedAt              time.Time     `json:"reportedAt,omitempty"`
	PriorHistoryUnavailable bool          `json:"priorHistoryUnavailable,omitempty"`
	StaleAfterRestart       bool          `json:"staleAfterRestart,omitempty"`
}

type ProviderCapability struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type ProviderStatus struct {
	Provider       string                        `json:"provider"`
	SetupState     string                        `json:"setupState"`
	Capabilities   map[string]ProviderCapability `json:"capabilities,omitempty"`
	Detail         string                        `json:"detail,omitempty"`
	Action         string                        `json:"action,omitempty"`
	AdapterVersion string                        `json:"adapterVersion,omitempty"`
}

type ObservationEnvelope struct {
	Version        int       `json:"version"`
	Provider       string    `json:"provider"`
	ProjectID      string    `json:"projectId"`
	Kind           string    `json:"kind"`
	SessionKey     string    `json:"sessionKey,omitempty"`
	AgentKey       string    `json:"agentKey,omitempty"`
	ParentAgentKey string    `json:"parentAgentKey,omitempty"`
	ToolKey        string    `json:"toolKey,omitempty"`
	ToolName       string    `json:"toolName,omitempty"`
	Outcome        string    `json:"outcome,omitempty"`
	ReportedAt     time.Time `json:"reportedAt,omitempty"`
}
type Project struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Root         string       `json:"root"`
	GitCommonDir string       `json:"gitCommonDir,omitempty"`
	DetectedAt   time.Time    `json:"detectedAt,omitempty"`
	Capabilities Capabilities `json:"capabilities,omitempty"`
}
type Runtime struct {
	ID              string    `json:"id"`
	ProjectID       string    `json:"projectId,omitempty"`
	Origin          string    `json:"origin"`
	ProtocolVersion int       `json:"protocolVersion"`
	StartedAt       time.Time `json:"startedAt"`
	State           string    `json:"state"`
}
type RuntimeRecord struct {
	RuntimeID       string    `json:"runtimeId"`
	ProjectID       string    `json:"projectId"`
	ProjectRoot     string    `json:"projectRoot"`
	Origin          string    `json:"origin"`
	PID             int       `json:"pid"`
	ProtocolVersion int       `json:"protocolVersion"`
	StartedAt       time.Time `json:"startedAt"`
	State           string    `json:"state"`
}
type DiagnosticCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
	Action string `json:"action,omitempty"`
}
type Snapshot struct {
	Sequence         uint64            `json:"sequence"`
	Project          Project           `json:"project"`
	Runtime          Runtime           `json:"runtime"`
	Diagnostics      []DiagnosticCheck `json:"diagnostics"`
	Agents           []Agent           `json:"agents,omitempty"`
	ProviderStatuses []ProviderStatus  `json:"providerStatuses,omitempty"`
}
type EventSource struct {
	Name           string `json:"name"`
	AdapterVersion string `json:"adapterVersion,omitempty"`
}
type Event struct {
	SchemaVersion int            `json:"schemaVersion"`
	EventID       string         `json:"eventId"`
	Sequence      uint64         `json:"sequence"`
	Timestamp     time.Time      `json:"timestamp"`
	ProjectID     string         `json:"projectId"`
	Source        EventSource    `json:"source"`
	Type          string         `json:"type"`
	Payload       map[string]any `json:"payload"`
}
type Health struct {
	Status           string            `json:"status"`
	Runtime          Runtime           `json:"runtime"`
	Project          Project           `json:"project"`
	Checks           []DiagnosticCheck `json:"checks"`
	Agents           []Agent           `json:"agents,omitempty"`
	ProviderStatuses []ProviderStatus  `json:"providerStatuses,omitempty"`
}
