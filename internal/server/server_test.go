package server

import (
	"github.com/craftlane/craftlane/internal/protocol"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRejectsInvalidHost(t *testing.T) {
	s := New(protocol.Project{ID: "p"}, protocol.Runtime{Origin: "http://127.0.0.1:7331"}, "x")
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7331/api/v1/health", nil)
	r.Host = "evil.test"
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatalf("got %d", w.Code)
	}
}
