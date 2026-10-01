package diagnostics

import (
	"github.com/craftlane/craftlane/internal/events"
	"github.com/craftlane/craftlane/internal/protocol"
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
