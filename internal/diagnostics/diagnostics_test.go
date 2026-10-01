package diagnostics

import "testing"

func TestFoundationFailureHasAction(t *testing.T) {
	for _, c := range Foundation(false, false, false) {
		if (c.Status == "fail" || c.Status == "warn") && c.Action == "" {
			t.Fatalf("missing action: %+v", c)
		}
	}
}
