package protocol

import "time"

type Capabilities map[string]string
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
	Sequence    uint64            `json:"sequence"`
	Project     Project           `json:"project"`
	Runtime     Runtime           `json:"runtime"`
	Diagnostics []DiagnosticCheck `json:"diagnostics"`
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
	Status  string            `json:"status"`
	Runtime Runtime           `json:"runtime"`
	Project Project           `json:"project"`
	Checks  []DiagnosticCheck `json:"checks"`
}
