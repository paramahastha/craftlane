package diagnostics

import (
	"github.com/craftlane/craftlane/internal/events"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/providers"
	"os"
	"os/exec"
	"path/filepath"
)

func Detect(projectRoot string, runtimeReady, dashboardReady bool) []protocol.DiagnosticCheck {
	_, gitErr := exec.LookPath("git")
	checks := Foundation(gitErr == nil, runtimeReady, dashboardReady)
	if _, err := os.Stat(filepath.Join(projectRoot, ".specify")); err == nil {
		checks[4] = protocol.DiagnosticCheck{Name: "speckit", Status: "pass", Detail: "Spec Kit project metadata detected"}
	}
	for _, executable := range []string{"claude", "codex"} {
		if _, err := exec.LookPath(executable); err == nil {
			checks[5] = protocol.DiagnosticCheck{Name: "provider", Status: "pass", Detail: executable + " provider executable detected"}
			break
		}
	}
	for i := range checks {
		checks[i].Detail = events.RedactText(checks[i].Detail)
		checks[i].Action = events.RedactText(checks[i].Action)
	}
	return checks
}

func DetectWithRegistry(projectRoot string, runtimeReady, dashboardReady bool, registry *providers.Registry) []protocol.DiagnosticCheck {
	checks := Detect(projectRoot, runtimeReady, dashboardReady)
	if registry == nil {
		return checks
	}
	for _, status := range registry.Statuses() {
		check := protocol.DiagnosticCheck{Name: "provider/" + status.Provider, Status: "pass", Detail: status.Detail, Action: status.Action}
		switch status.SetupState {
		case "unavailable":
			check.Status = "warn"
			if check.Action == "" {
				check.Action = "Install or enable the provider"
			}
		case "not_configured", "needs_trust", "degraded":
			check.Status = "warn"
			if check.Action == "" {
				check.Action = "Run craftlane observe setup " + status.Provider
			}
		}
		if check.Detail == "" {
			check.Detail = "Provider status: " + status.SetupState
		}
		checks = append(checks, check)
	}
	return checks
}
