package integration

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestEventsHTTPContractAndCursorValidation(t *testing.T) {
	runtime, _ := startRuntime(t)
	request, _ := http.NewRequest(http.MethodGet, runtime.Record.Origin+"/api/v1/events?after=0", nil)
	request.Host = request.URL.Host
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var body struct {
		Status string            `json:"status"`
		Cursor uint64            `json:"cursor"`
		Events []json.RawMessage `json:"events"`
	}
	if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 || body.Status != "ok" || body.Cursor < 2 || len(body.Events) < 2 {
		t.Fatalf("status=%d body=%#v", response.StatusCode, body)
	}
	if status := getJSON(t, runtime.Record.Origin, "/api/v1/events?after=bad", nil); status != 400 {
		t.Fatalf("invalid cursor status=%d", status)
	}
}

func TestEventSchemaDeclaresRequiredEnvelopeFields(t *testing.T) {
	data, err := os.ReadFile("../../protocol/schema/events.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"schemaVersion", "eventId", "sequence", "timestamp", "projectId", "source", "type", "payload"} {
		if !strings.Contains(string(data), `"`+field+`"`) {
			t.Errorf("schema missing %s", field)
		}
	}
}
