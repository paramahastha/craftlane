package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/providers"
	"github.com/craftlane/craftlane/internal/providers/claude"
	"github.com/craftlane/craftlane/internal/providers/codex"
	"github.com/craftlane/craftlane/internal/providers/opencode"
)

func NewProviderRegistry(root string) *providers.Registry {
	return providers.NewRegistry(claude.New(root), codex.New(root), opencode.New(root))
}

func ObserveRegistry(p protocol.Project) *providers.Registry { return NewProviderRegistry(p.Root) }

func ObserveSetup(p protocol.Project, provider string, confirm bool) (providers.SetupResult, error) {
	adapter, ok := ObserveRegistry(p).Get(provider)
	if !ok {
		return providers.SetupResult{}, fmt.Errorf("unsupported provider %q", provider)
	}
	preview, err := adapter.Preview()
	if err != nil {
		return providers.SetupResult{}, err
	}
	if !confirm {
		return providers.SetupResult{Preview: preview}, nil
	}
	return adapter.Apply()
}

func ObserveRemove(p protocol.Project, provider string, confirm bool) (providers.SetupResult, error) {
	registry := ObserveRegistry(p)
	adapter, ok := registry.Get(provider)
	if !ok {
		return providers.SetupResult{}, fmt.Errorf("unsupported provider %q", provider)
	}
	preview, err := adapter.Preview()
	if err != nil {
		return providers.SetupResult{}, err
	}
	if !confirm {
		return providers.SetupResult{Preview: preview}, nil
	}
	result, err := adapter.Remove()
	if err != nil {
		return result, err
	}
	needed := false
	for _, id := range registry.IDs() {
		if other, ok := registry.Get(id); ok {
			if otherPreview, previewErr := other.Preview(); previewErr == nil && otherPreview.AlreadyApplied {
				needed = true
			}
		}
	}
	if !needed {
		if store, storeErr := StorageForProject(p); storeErr == nil {
			_ = store.ClearObservationSecret()
		}
	}
	return result, nil
}

func ObserveStatus(p protocol.Project) []protocol.ProviderStatus {
	return ObserveRegistry(p).Statuses()
}

func ObserveBridge(ctx context.Context, p protocol.Project, provider string, input io.Reader) error {
	if !providers.Allowed[provider] {
		return fmt.Errorf("unsupported provider %q", provider)
	}
	var raw map[string]any
	decoder := json.NewDecoder(io.LimitReader(input, providers.MaxBridgeInput))
	if err := decoder.Decode(&raw); err != nil {
		return fmt.Errorf("invalid bridge input")
	}
	if raw["projectId"] == nil {
		raw["projectId"] = p.ID
	}
	envelope, err := providers.ProjectInput(provider, raw)
	if err != nil {
		return err
	}
	store, err := StorageForProject(p)
	if err != nil {
		return err
	}
	record, live, err := RuntimeStatus(ctx, p)
	if err != nil {
		return err
	}
	if !live {
		return fmt.Errorf("no live runtime; run craftlane start")
	}
	token, err := store.LoadObservationSecret()
	if err != nil {
		return fmt.Errorf("observation secret unavailable")
	}
	return providers.Send(ctx, record.Origin, token, provider, envelope)
}

func PrintSetupResult(w io.Writer, result providers.SetupResult) {
	for _, change := range result.Preview.Changes {
		fmt.Fprintf(w, "%s: %s — %s\n", change.Action, change.Path, change.Detail)
	}
	for _, warning := range result.Preview.Warnings {
		fmt.Fprintf(w, "warning: %s\n", warning)
	}
	if result.Preview.RequiresTrust {
		fmt.Fprintln(w, "trust: provider hook review is required")
	}
	if result.Changed {
		fmt.Fprintln(w, "changed: yes")
	}
}
