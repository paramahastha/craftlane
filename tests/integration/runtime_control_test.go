package integration

import (
	"context"
	"github.com/craftlane/craftlane/internal/app"
	"net/http"
	"testing"
)

func TestRuntimeControlRejectsMissingInvalidAndCrossOriginTokens(t *testing.T) {
	runtime, project := startRuntime(t)
	oldToken := runtime.Token
	tests := []struct {
		token, origin string
		want          int
	}{
		{"", "", http.StatusUnauthorized},
		{"wrong", "", http.StatusUnauthorized},
		{runtime.Token, "http://evil.test", http.StatusForbidden},
	}
	for _, test := range tests {
		req, _ := http.NewRequest(http.MethodPost, runtime.Record.Origin+"/api/v1/runtime/stop", nil)
		req.Host = req.URL.Host
		req.Header.Set("X-Craftlane-Token", test.token)
		if test.origin != "" {
			req.Header.Set("Origin", test.origin)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != test.want {
			t.Fatalf("token/origin status=%d want=%d", resp.StatusCode, test.want)
		}
	}
	if err := runtime.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.Store.LoadControl(); err == nil {
		t.Fatal("shutdown token remains readable")
	}
	replacement, err := app.Start(context.Background(), project, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer replacement.Shutdown(context.Background())
	req, _ := http.NewRequest(http.MethodPost, replacement.Record.Origin+"/api/v1/runtime/stop", nil)
	req.Host = req.URL.Host
	req.Header.Set("X-Craftlane-Token", oldToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("obsolete token status=%d", resp.StatusCode)
	}
}
