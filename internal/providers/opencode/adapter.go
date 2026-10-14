package opencode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/providers"
)

type Adapter struct{ Root string }

func New(root string) *Adapter     { return &Adapter{Root: root} }
func (a *Adapter) ID() string      { return protocol.ProviderOpenCode }
func (a *Adapter) Version() string { return "opencode-plugin-v1" }
func (a *Adapter) config() string {
	return filepath.Join(a.Root, ".opencode", "plugins", "craftlane-observability.ts")
}
func (a *Adapter) Status() protocol.ProviderStatus {
	state, detail := providers.ExecutableStatusForTests("opencode")
	setup := "not_configured"
	if state == "unavailable" {
		setup = "unavailable"
	} else if preview, err := a.Preview(); err != nil {
		setup = "degraded"
	} else if preview.AlreadyApplied {
		setup = "ready"
	}
	return protocol.ProviderStatus{Provider: a.ID(), SetupState: setup, AdapterVersion: a.Version(), Detail: detail, Action: "Install a compatible OpenCode plugin API", Capabilities: map[string]protocol.ProviderCapability{"sessions": {Status: "available"}, "subagents": {Status: "available"}, "tools": {Status: "available"}, "permissions": {Status: "available"}, "failures": {Status: "available"}, "terminal": {Status: "unavailable", Reason: "OpenCode status has no durable successful terminal signal"}, "idle": {Status: "partial", Reason: "Busy/idle status is provider-reported"}}}
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
	b, err := os.ReadFile(a.config())
	if err != nil && !os.IsNotExist(err) {
		return providers.SetupPreview{}, err
	}
	preview := providers.SetupPreview{Provider: a.ID(), Changes: []providers.SetupChange{{Path: a.config(), Action: "create", Detail: "Create one version-pinned Craftlane-owned plugin"}}, Warnings: []string{"The generated project plugin may be Git-tracked; Craftlane does not edit ignore rules"}}
	if err == nil {
		if !strings.Contains(string(b), "Craftlane-owned observability plugin v1") {
			return providers.SetupPreview{}, fmt.Errorf("existing OpenCode plugin is not Craftlane-owned")
		}
		preview.AlreadyApplied = true
		preview.Changes[0].Action = "retain"
		preview.Changes[0].Detail = "Craftlane observation plugin already present"
	}
	return preview, nil
}
func (a *Adapter) Apply() (providers.SetupResult, error) {
	preview, err := a.Preview()
	if err != nil {
		return providers.SetupResult{}, err
	}
	if preview.AlreadyApplied {
		return providers.SetupResult{Preview: preview}, nil
	}
	if err := os.MkdirAll(filepath.Dir(a.config()), 0700); err != nil {
		return providers.SetupResult{}, err
	}
	source := "// Craftlane-owned observability plugin v1\nexport const craftlaneProvider = 'opencode';\nexport const craftlaneEventKinds = ['session.started', 'tool.started', 'tool.completed', 'permission.requested'];\n"
	if err := os.WriteFile(a.config(), []byte(source), 0600); err != nil {
		return providers.SetupResult{}, err
	}
	return providers.SetupResult{Preview: preview, Changed: true}, nil
}
func (a *Adapter) Remove() (providers.SetupResult, error) {
	preview, err := a.Preview()
	if err != nil {
		return providers.SetupResult{}, err
	}
	if preview.AlreadyApplied == false {
		return providers.SetupResult{Preview: preview}, nil
	}
	err = os.Remove(a.config())
	if os.IsNotExist(err) {
		err = nil
	}
	return providers.SetupResult{Preview: preview, Changed: err == nil}, err
}
