package integration

import (
	"github.com/craftlane/craftlane/internal/protocol"
	"testing"
)

func TestProjectAndHealthContracts(t *testing.T) {
	runtime, project := startRuntime(t)
	var health protocol.Health
	if status := getJSON(t, runtime.Record.Origin, "/api/v1/health", &health); status != 200 {
		t.Fatalf("health status = %d", status)
	}
	if health.Runtime.ID != runtime.Record.RuntimeID || health.Project.ID != project.ID || len(health.Checks) < 5 {
		t.Fatalf("health = %#v", health)
	}
	var got protocol.Project
	if status := getJSON(t, runtime.Record.Origin, "/api/v1/project", &got); status != 200 || got.Root != project.Root {
		t.Fatalf("project status=%d value=%#v", status, got)
	}
}
