package codex

import (
	"fmt"
	"path/filepath"

	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/providers"
)

type Adapter struct{ Root string }

func New(root string) *Adapter     { return &Adapter{Root: root} }
func (a *Adapter) ID() string      { return protocol.ProviderCodex }
func (a *Adapter) Version() string { return "codex-hooks-v1" }
func (a *Adapter) config() string  { return filepath.Join(a.Root, ".codex", "hooks.json") }
func (a *Adapter) Status() protocol.ProviderStatus {
	state, detail := providers.ExecutableStatusForTests("codex")
	setup := "not_configured"
	if state == "unavailable" {
		setup = "unavailable"
	} else if preview, err := a.Preview(); err != nil {
		setup = "degraded"
	} else if preview.AlreadyApplied {
		setup = "needs_trust"
	}
	return protocol.ProviderStatus{Provider: a.ID(), SetupState: setup, AdapterVersion: a.Version(), Detail: detail, Action: "Review and trust the project hook if Codex requests it", Capabilities: map[string]protocol.ProviderCapability{"sessions": {Status: "available"}, "subagents": {Status: "unavailable", Reason: "Codex hooks do not expose subagent lifecycle"}, "tools": {Status: "partial", Reason: "Only supported hook categories are visible"}, "permissions": {Status: "available"}, "failures": {Status: "available"}, "terminal": {Status: "partial", Reason: "SessionEnd may be delayed"}, "idle": {Status: "unavailable", Reason: "No durable idle hook"}}}
}
func (a *Adapter) Normalize(e protocol.ObservationEnvelope) (protocol.ObservationEnvelope, error) {
	if e.Provider != a.ID() {
		return protocol.ObservationEnvelope{}, fmt.Errorf("provider mismatch")
	}
	if err := providers.ValidateEnvelope(e); err != nil {
		return protocol.ObservationEnvelope{}, err
	}
	if e.Kind == "turn.completed" {
		e.Kind = "turn.updated"
	}
	return e, nil
}
func (a *Adapter) Preview() (providers.SetupPreview, error) {
	return providers.PreviewFile(a.config(), a.ID(), "Codex project hooks may be Git-tracked; Craftlane does not edit ignore rules", true)
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
