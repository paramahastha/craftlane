package diagnostics

import (
	"fmt"
	"github.com/craftlane/craftlane/internal/protocol"
)

func Render(checks []protocol.DiagnosticCheck) string {
	out := ""
	for _, c := range checks {
		out += fmt.Sprintf("%s: %s — %s", c.Name, c.Status, c.Detail)
		if c.Action != "" {
			out += " (next: " + c.Action + ")"
		}
		out += "\n"
	}
	return out
}
