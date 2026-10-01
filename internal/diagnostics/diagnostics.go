package diagnostics

import "github.com/craftlane/craftlane/internal/protocol"

func Foundation(gitAvailable, runtimeReady, dashboardReady bool) []protocol.DiagnosticCheck {
	checks := []protocol.DiagnosticCheck{{Name: "git", Status: "pass", Detail: "Git is available"}, {Name: "project-detection", Status: "pass", Detail: "Git project was detected"}, {Name: "runtime-record", Status: "pass", Detail: "Runtime state is readable"}, {Name: "embedded-dashboard", Status: "pass", Detail: "Embedded dashboard is available"}}
	if !gitAvailable {
		checks[0] = protocol.DiagnosticCheck{Name: "git", Status: "fail", Detail: "Git is unavailable", Action: "Install Git and retry"}
	}
	if !runtimeReady {
		checks[2] = protocol.DiagnosticCheck{Name: "runtime-record", Status: "fail", Detail: "Runtime is not ready", Action: "Run craftlane start"}
	}
	if !dashboardReady {
		checks[3] = protocol.DiagnosticCheck{Name: "embedded-dashboard", Status: "fail", Detail: "Dashboard asset is unavailable", Action: "Run make web-build"}
	}
	checks = append(checks, protocol.DiagnosticCheck{Name: "speckit", Status: "warn", Detail: "Spec Kit is optional and unavailable", Action: "Install Spec Kit if needed"}, protocol.DiagnosticCheck{Name: "provider", Status: "warn", Detail: "Provider integrations are optional", Action: "Configure a provider if needed"})
	return checks
}
