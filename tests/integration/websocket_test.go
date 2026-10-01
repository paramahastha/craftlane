package integration

import (
	"context"
	"encoding/json"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/craftlane/craftlane/internal/protocol"
	"net/http"
	"strings"
	"testing"
	"time"
)

type wsMessage struct {
	Type   string `json:"type"`
	Cursor uint64 `json:"cursor"`
	Event  struct {
		EventID  string `json:"eventId"`
		Sequence uint64 `json:"sequence"`
	} `json:"event"`
}

func TestWebSocketResyncsCursorOutsideRetention(t *testing.T) {
	runtime, _ := startRuntime(t)
	runtime.Reducer.Events = []protocol.Event{{Sequence: 5, EventID: "retained", ProjectID: runtime.Project.ID}}
	runtime.Reducer.Sequence = 5
	conn, first := dialEvents(t, runtime.Record.Origin, 1)
	defer conn.CloseNow()
	if first.Type != "resync_required" || first.Cursor != 5 {
		t.Fatalf("first=%#v", first)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var snapshot wsMessage
	if err := wsjson.Read(ctx, conn, &snapshot); err != nil || snapshot.Type != "snapshot" || snapshot.Cursor != 5 {
		t.Fatalf("snapshot=%#v err=%v", snapshot, err)
	}
}

func dialEvents(t *testing.T, origin string, last uint64) (*websocket.Conn, wsMessage) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(cancel)
	url := "ws" + strings.TrimPrefix(origin, "http") + "/ws"
	conn, response, err := websocket.Dial(ctx, url, &websocket.DialOptions{HTTPHeader: http.Header{"Origin": []string{origin}}, Subprotocols: []string{"craftlane.events.v1"}})
	if err != nil {
		if response != nil {
			t.Fatalf("dial status=%d: %v", response.StatusCode, err)
		}
		t.Fatal(err)
	}
	if conn.Subprotocol() != "craftlane.events.v1" {
		t.Fatalf("subprotocol=%q", conn.Subprotocol())
	}
	if err = wsjson.Write(ctx, conn, map[string]any{"type": "subscribe", "protocolVersion": 1, "lastSequence": last}); err != nil {
		t.Fatal(err)
	}
	var message wsMessage
	if err = wsjson.Read(ctx, conn, &message); err != nil {
		t.Fatal(err)
	}
	return conn, message
}

func TestWebSocketOriginProtocolReplayAndReconnects(t *testing.T) {
	runtime, _ := startRuntime(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	url := "ws" + strings.TrimPrefix(runtime.Record.Origin, "http") + "/ws"
	if conn, response, err := websocket.Dial(ctx, url, &websocket.DialOptions{HTTPHeader: http.Header{"Origin": []string{"http://evil.test"}}, Subprotocols: []string{"craftlane.events.v1"}}); err == nil || response.StatusCode != http.StatusForbidden {
		if conn != nil {
			conn.CloseNow()
		}
		t.Fatalf("cross-origin response=%v err=%v", response, err)
	}
	if conn, response, err := websocket.Dial(ctx, url, &websocket.DialOptions{HTTPHeader: http.Header{"Origin": []string{runtime.Record.Origin}}}); err == nil || response.StatusCode != http.StatusBadRequest {
		if conn != nil {
			conn.CloseNow()
		}
		t.Fatalf("missing protocol response=%v err=%v", response, err)
	}
	conn, first := dialEvents(t, runtime.Record.Origin, 0)
	if first.Type != "snapshot" {
		t.Fatalf("first message=%#v", first)
	}
	cursor := first.Cursor
	liveEvent, err := runtime.Reducer.Append("project.changed", map[string]any{"path": "live.txt"})
	if err != nil {
		t.Fatal(err)
	}
	var live wsMessage
	if err = wsjson.Read(ctx, conn, &live); err != nil || live.Type != "event" || live.Event.EventID != liveEvent.EventID {
		t.Fatalf("live event=%#v err=%v", live, err)
	}
	cursor = live.Event.Sequence
	conn.CloseNow()
	event, err := runtime.Reducer.Append("project.changed", map[string]any{"path": "README.md"})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		conn, replay := dialEvents(t, runtime.Record.Origin, cursor)
		if replay.Type != "event" || replay.Event.Sequence != event.Sequence || replay.Event.EventID != event.EventID {
			raw, _ := json.Marshal(replay)
			t.Fatalf("reconnect %d: %s", i, raw)
		}
		conn.CloseNow()
	}
}
