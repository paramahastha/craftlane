package integration

import (
	"context"
	"github.com/craftlane/craftlane/internal/app"
	"testing"
)

func TestDoctorReportsRequiredAndOptionalChecks(t *testing.T) {
	_, project := startRuntime(t)
	checks, ok := app.Doctor(context.Background(), project)
	if !ok {
		t.Fatal("required diagnostics unexpectedly failed")
	}
	want := map[string]bool{"git": false, "runtime-record": false, "embedded-dashboard": false, "speckit": false, "provider": false}
	for _, check := range checks {
		if _, exists := want[check.Name]; exists {
			want[check.Name] = true
		}
		if (check.Status == "warn" || check.Status == "fail") && check.Action == "" {
			t.Fatalf("missing action: %#v", check)
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing %s diagnostic", name)
		}
	}
}
