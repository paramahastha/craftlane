package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/craftlane/craftlane/internal/protocol"
)

const MaxBridgeInput = 16 << 10

func ProjectInput(provider string, input map[string]any) (protocol.ObservationEnvelope, error) {
	if !Allowed[provider] {
		return protocol.ObservationEnvelope{}, fmt.Errorf("unsupported provider %q", provider)
	}
	get := func(key string) string { value, _ := input[key].(string); return strings.TrimSpace(value) }
	envelope := protocol.ObservationEnvelope{
		Version: 1, Provider: provider, ProjectID: get("projectId"), Kind: get("kind"),
		SessionKey: get("sessionKey"), AgentKey: get("agentKey"), ParentAgentKey: get("parentAgentKey"),
		ToolKey: get("toolKey"), ToolName: get("toolName"), Outcome: get("outcome"),
	}
	if reported := get("reportedAt"); reported != "" {
		parsed, err := time.Parse(time.RFC3339, reported)
		if err != nil {
			return protocol.ObservationEnvelope{}, fmt.Errorf("invalid reportedAt")
		}
		envelope.ReportedAt = parsed
	}
	if err := ValidateEnvelope(envelope); err != nil {
		return protocol.ObservationEnvelope{}, err
	}
	return envelope, nil
}

func ValidateEnvelope(e protocol.ObservationEnvelope) error {
	if e.Version != 1 || !Allowed[e.Provider] || e.ProjectID == "" || e.Kind == "" {
		return fmt.Errorf("invalid observation envelope")
	}
	if e.AgentKey == "" && e.SessionKey == "" {
		return fmt.Errorf("observation requires a session or agent key")
	}
	if len(e.ProjectID) > 256 || len(e.Kind) > 64 || len(e.SessionKey) > 256 || len(e.AgentKey) > 256 || len(e.ToolName) > 128 {
		return fmt.Errorf("observation envelope is too large")
	}
	return nil
}

func Send(ctx context.Context, origin, token, provider string, envelope protocol.ObservationEnvelope) error {
	if err := ValidateEnvelope(envelope); err != nil {
		return err
	}
	b, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(origin, "/")+"/api/v1/observe/"+provider, bytes.NewReader(b))
	if err != nil {
		return err
	}
	request.Host = request.URL.Host
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Craftlane-Observation-Token", token)
	client := &http.Client{Timeout: 750 * time.Millisecond}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 512))
		return fmt.Errorf("observation callback returned %s", response.Status)
	}
	return nil
}
