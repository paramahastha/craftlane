package integration

import (
	"context"
	"encoding/json"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/craftlane/craftlane/internal/protocol"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSecretsNeverReachPersistenceHTTPOrDiagnostics(t *testing.T) {
	runtime, _ := startRuntime(t)
	secret := "sk-super-secret-value"
	event, err := runtime.Reducer.Append("project.changed", map[string]any{"authorization": secret, "message": "Bearer " + secret})
	if err != nil {
		t.Fatal(err)
	}
	if err = runtime.Reducer.SetDiagnostics([]protocol.DiagnosticCheck{{Name: "safe", Status: "warn", Detail: secret, Action: "remove " + secret}}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"events.jsonl", "snapshot.json"} {
		data, err := os.ReadFile(runtime.Store.Dir + "/" + name)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), secret) {
			t.Fatalf("secret reached %s", name)
		}
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, runtime.Record.Origin+"/api/v1/events", nil)
	req.Host = req.URL.Host
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if strings.Contains(string(body), secret) {
		t.Fatal("secret reached HTTP response")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	wsURL := "ws" + strings.TrimPrefix(runtime.Record.Origin, "http") + "/ws"
	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{HTTPHeader: http.Header{"Origin": []string{runtime.Record.Origin}}, Subprotocols: []string{"craftlane.events.v1"}})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.CloseNow()
	if err = wsjson.Write(ctx, conn, map[string]any{"type": "subscribe", "protocolVersion": 1, "lastSequence": event.Sequence - 1}); err != nil {
		t.Fatal(err)
	}
	var message map[string]any
	if err = wsjson.Read(ctx, conn, &message); err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(message)
	if strings.Contains(string(encoded), secret) {
		t.Fatal("secret reached WebSocket response")
	}
}
