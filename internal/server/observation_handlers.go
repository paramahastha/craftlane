package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/providers"
	"github.com/craftlane/craftlane/internal/security"
)

const maxObservationBody = 16 << 10

func (s *Server) observe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Origin") != "" && !security.Origin(r.Header.Get("Origin"), s.ExpectedOrigin) {
		http.Error(w, "forbidden origin", http.StatusForbidden)
		return
	}
	if !security.Token(r.Header.Get("X-Craftlane-Observation-Token"), s.ObservationToken) {
		http.Error(w, "invalid observation token", http.StatusUnauthorized)
		return
	}
	provider := strings.TrimPrefix(r.URL.Path, "/api/v1/observe/")
	if provider == "" || strings.Contains(provider, "/") || !providers.Allowed[provider] {
		http.Error(w, "unsupported provider", http.StatusNotFound)
		return
	}
	if s.Registry == nil || s.Reducer == nil {
		http.Error(w, "observation unavailable", http.StatusServiceUnavailable)
		return
	}
	var envelope protocol.ObservationEnvelope
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxObservationBody))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		http.Error(w, "invalid observation envelope", http.StatusBadRequest)
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "invalid observation envelope", http.StatusBadRequest)
		return
	}
	if envelope.Provider != provider {
		http.Error(w, "provider mismatch", http.StatusBadRequest)
		return
	}
	normalized, adapter, err := s.Registry.Normalize(envelope)
	if err != nil || normalized.ProjectID != s.Project.ID {
		http.Error(w, safeObservationError(err), http.StatusBadRequest)
		return
	}
	_, changed, err := s.Reducer.ApplyObservation(normalized, protocol.EventSource{Name: adapter.ID(), AdapterVersion: adapter.Version()})
	if err != nil {
		http.Error(w, safeObservationError(err), http.StatusBadRequest)
		return
	}
	jsonWrite(w, map[string]any{"accepted": true, "changed": changed})
}

func safeObservationError(err error) string {
	if err == nil {
		return "invalid observation"
	}
	return fmt.Sprintf("invalid observation: %s", err.Error())
}
