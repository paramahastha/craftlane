package claude

import (
	"fmt"
	"path/filepath"

	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/providers"
)

type Adapter struct{ Root string }

func New(root string) *Adapter     { return &Adapter{Root: root} }
func (a *Adapter) ID() string      { return protocol.ProviderClaude }
func (a *Adapter) Version() string { return "claude-hooks-v1" }
func (a *Adapter) config() string  { return filepath.Join(a.Root, ".claude", "settings.local.json") }
func (a *Adapter) Status() protocol.ProviderStatus {
	state, detail := providers.ExecutableStatusForTests("claude")
	setup := "not_configured"
	if state != "unavailable" {
		if preview, err := a.Preview(); err != nil {
			setup = "degraded"
		} else if preview.AlreadyApplied {
			setup = "ready"
		}
	}
	if state == "unavailable" {
		setup = "unavailable"
	}
	return protocol.ProviderStatus{Provider: a.ID(), SetupState: setup, AdapterVersion: a.Version(), Detail: detail, Capabilities: capabilities()}
}
func capabilities() map[string]protocol.ProviderCapability {
	return map[string]protocol.ProviderCapability{"sessions": {Status: "available"}, "subagents": {Status: "available"}, "tools": {Status: "available"}, "permissions": {Status: "available"}, "failures": {Status: "available"}, "terminal": {Status: "available"}, "idle": {Status: "unavailable", Reason: "Claude hooks do not report durable idle state"}}
}
func (a *Adapter) Normalize(e protocol.ObservationEnvelope) (protocol.ObservationEnvelope, error) {
	if e.Provider != a.ID() {
		return protocol.ObservationEnvelope{}, fmt.Errorf("provider mismatch")
	}
	if err := providers.ValidateEnvelope(e); err != nil {
		return protocol.ObservationEnvelope{}, err
	}
	return e, nil
}
func (a *Adapter) Preview() (providers.SetupPreview, error) {
	return providers.PreviewFile(a.config(), a.ID(), "Settings are project-local and preserve unrelated hooks", false)
}
func (a *Adapter) Apply() (providers.SetupResult, error) {
	preview, err := a.Preview()
	if err != nil {
		return providers.SetupResult{}, err
	}
	changed, err := providers.ApplyFile(a.config(), a.ID())
	return providers.SetupResult{Preview: preview, Changed: changed}, err
}
func (a *Adapter) Remove() (providers.SetupResult, error) {
	preview, err := a.Preview()
	if err != nil {
		return providers.SetupResult{}, err
	}
	changed, err := providers.RemoveFile(a.config(), a.ID())
	return providers.SetupResult{Preview: preview, Changed: changed}, err
}
