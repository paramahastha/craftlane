package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/craftlane/craftlane/internal/events"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/providers"
	"github.com/craftlane/craftlane/internal/providers/claude"
	"github.com/craftlane/craftlane/internal/storage"
)

func TestObservationCallbackGuardsAndMutation(t *testing.T) {
	p := protocol.Project{ID: "p", Root: t.TempDir()}
	r := events.NewReducer(storage.New(t.TempDir()), p, protocol.Runtime{})
	s := New(p, protocol.Runtime{Origin: "http://127.0.0.1:7331"}, "control")
	s.ObservationToken = "observe"
	s.Registry = providers.NewRegistry(claude.New(p.Root))
	s.Reducer = r
	h := s.Handler()
	body := []byte(`{"version":1,"provider":"claude","projectId":"p","kind":"session.started","sessionKey":"s","agentKey":"a"}`)
	request := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:7331/api/v1/observe/claude", bytes.NewReader(body))
	request.Host = "127.0.0.1:7331"
	request.Header.Set("X-Craftlane-Observation-Token", "observe")
	response := httptest.NewRecorder()
	h.ServeHTTP(response, request)
	if response.Code != http.StatusOK || len(r.Agents()) != 1 {
		t.Fatalf("accepted callback: status=%d agents=%d", response.Code, len(r.Agents()))
	}
	request = httptest.NewRequest(http.MethodPost, "http://127.0.0.1:7331/api/v1/observe/claude", bytes.NewReader(body))
	request.Host = "evil.test"
	request.Header.Set("X-Craftlane-Observation-Token", "observe")
	response = httptest.NewRecorder()
	h.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("wrong host status=%d", response.Code)
	}
}
