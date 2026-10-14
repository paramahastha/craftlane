package providers

import (
	"fmt"
	"sort"

	"github.com/craftlane/craftlane/internal/protocol"
)

var Allowed = map[string]bool{
	protocol.ProviderClaude:   true,
	protocol.ProviderCodex:    true,
	protocol.ProviderOpenCode: true,
}

type SetupChange struct {
	Path   string `json:"path"`
	Action string `json:"action"`
	Detail string `json:"detail"`
}

type SetupPreview struct {
	Provider       string        `json:"provider"`
	Changes        []SetupChange `json:"changes"`
	Warnings       []string      `json:"warnings,omitempty"`
	RequiresTrust  bool          `json:"requiresTrust,omitempty"`
	AlreadyApplied bool          `json:"alreadyApplied,omitempty"`
}

type SetupResult struct {
	Preview SetupPreview `json:"preview"`
	Changed bool         `json:"changed"`
}

type Adapter interface {
	ID() string
	Version() string
	Status() protocol.ProviderStatus
	Normalize(protocol.ObservationEnvelope) (protocol.ObservationEnvelope, error)
	Preview() (SetupPreview, error)
	Apply() (SetupResult, error)
	Remove() (SetupResult, error)
}

type Registry struct{ adapters map[string]Adapter }

func NewRegistry(adapters ...Adapter) *Registry {
	r := &Registry{adapters: make(map[string]Adapter)}
	for _, adapter := range adapters {
		if adapter != nil {
			r.adapters[adapter.ID()] = adapter
		}
	}
	return r
}

func (r *Registry) IDs() []string {
	ids := make([]string, 0, len(r.adapters))
	for id := range r.adapters {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func (r *Registry) Get(id string) (Adapter, bool) {
	adapter, ok := r.adapters[id]
	return adapter, ok && Allowed[id]
}

func (r *Registry) Statuses() []protocol.ProviderStatus {
	statuses := make([]protocol.ProviderStatus, 0, len(r.adapters))
	for _, id := range r.IDs() {
		statuses = append(statuses, r.adapters[id].Status())
	}
	return statuses
}

func (r *Registry) Normalize(envelope protocol.ObservationEnvelope) (protocol.ObservationEnvelope, Adapter, error) {
	if !Allowed[envelope.Provider] {
		return protocol.ObservationEnvelope{}, nil, fmt.Errorf("unsupported provider %q", envelope.Provider)
	}
	adapter, ok := r.Get(envelope.Provider)
	if !ok {
		return protocol.ObservationEnvelope{}, nil, fmt.Errorf("provider %q is unavailable", envelope.Provider)
	}
	normalized, err := adapter.Normalize(envelope)
	return normalized, adapter, err
}
