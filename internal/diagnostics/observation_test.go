package diagnostics

import (
	"testing"

	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/providers"
)

func TestProviderDiagnosticsStayIndependentAndActionable(t *testing.T) {
	r := providers.NewRegistry()
	checks := DetectWithRegistry(t.TempDir(), true, true, r)
	if len(checks) < 5 {
		t.Fatalf("checks = %#v", checks)
	}
	for _, check := range checks {
		if (check.Status == "warn" || check.Status == "fail") && check.Action == "" {
			t.Fatalf("missing action: %#v", check)
		}
	}
	_ = protocol.ProviderStatus{}
}
