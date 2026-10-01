package server

import (
	"github.com/craftlane/craftlane/internal/protocol"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmbeddedAssetsAndSPAFallback(t *testing.T) {
	if !AssetReady() {
		t.Fatal("embedded dashboard is missing")
	}
	server := New(protocol.Project{ID: "p"}, protocol.Runtime{Origin: "http://127.0.0.1:7331"}, "token")
	for _, path := range []string{"/", "/workspace/deep-link"} {
		request := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7331"+path, nil)
		request.Host = "127.0.0.1:7331"
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "root") {
			t.Fatalf("path %s: status=%d body=%q", path, response.Code, response.Body.String())
		}
	}
}
