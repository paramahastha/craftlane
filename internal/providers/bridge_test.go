package providers

import "testing"

func TestProjectInputAllowlistsSafeFields(t *testing.T) {
	e, err := ProjectInput("claude", map[string]any{"projectId": "p", "kind": "tool.started", "sessionKey": "s", "agentKey": "a", "toolName": "apply_patch", "prompt": "do not retain", "toolInput": "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if e.ToolName != "apply_patch" || e.ProjectID != "p" {
		t.Fatalf("envelope = %#v", e)
	}
}

func TestProjectInputRejectsMalformedAndUnknownProvider(t *testing.T) {
	if _, err := ProjectInput("other", map[string]any{}); err == nil {
		t.Fatal("unknown provider accepted")
	}
	if _, err := ProjectInput("codex", map[string]any{"projectId": "p", "kind": "tool.started"}); err == nil {
		t.Fatal("missing identity accepted")
	}
}
