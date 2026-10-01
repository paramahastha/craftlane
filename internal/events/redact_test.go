package events

import "testing"

func TestRedactRecursive(t *testing.T) {
	v := Redact(map[string]any{"token": "secret-value", "nested": []any{map[string]any{"password": "pw"}}, "text": "Bearer abc"}).(map[string]any)
	if v["token"] != "[REDACTED]" || v["nested"].([]any)[0].(map[string]any)["password"] != "[REDACTED]" {
		t.Fatalf("not redacted: %#v", v)
	}
}
