package providers

import (
	"testing"

	"github.com/craftlane/craftlane/internal/protocol"
)

type fakeAdapter struct{ id string }

func (f fakeAdapter) ID() string                      { return f.id }
func (f fakeAdapter) Version() string                 { return "test" }
func (f fakeAdapter) Status() protocol.ProviderStatus { return protocol.ProviderStatus{Provider: f.id} }
func (f fakeAdapter) Normalize(e protocol.ObservationEnvelope) (protocol.ObservationEnvelope, error) {
	return e, nil
}
func (f fakeAdapter) Preview() (SetupPreview, error) { return SetupPreview{Provider: f.id}, nil }
func (f fakeAdapter) Apply() (SetupResult, error)    { return SetupResult{}, nil }
func (f fakeAdapter) Remove() (SetupResult, error)   { return SetupResult{}, nil }

func TestRegistryOnlyReturnsAllowlistedAdapters(t *testing.T) {
	r := NewRegistry(fakeAdapter{id: protocol.ProviderClaude}, fakeAdapter{id: "other"})
	if _, ok := r.Get(protocol.ProviderClaude); !ok {
		t.Fatal("claude missing")
	}
	if _, ok := r.Get("other"); ok {
		t.Fatal("unknown adapter returned")
	}
	if _, _, err := r.Normalize(protocol.ObservationEnvelope{Provider: "other"}); err == nil {
		t.Fatal("unknown provider normalized")
	}
}
