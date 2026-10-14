package opencode

import (
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/providers"
)

func Project(input map[string]any) (protocol.ObservationEnvelope, error) {
	return providers.ProjectInput("opencode", input)
}
