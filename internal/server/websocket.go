package server

import (
	"context"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/craftlane/craftlane/internal/security"
	"net/http"
	"strings"
	"time"
)

const eventSubprotocol = "craftlane.events.v1"

type subscribeMessage struct {
	Type            string `json:"type"`
	ProtocolVersion int    `json:"protocolVersion"`
	LastSequence    uint64 `json:"lastSequence"`
}

func (s *Server) websocket(w http.ResponseWriter, r *http.Request) {
	if s.Reducer == nil {
		http.Error(w, "events unavailable", http.StatusServiceUnavailable)
		return
	}
	if !security.Origin(r.Header.Get("Origin"), s.ExpectedOrigin) {
		http.Error(w, "invalid origin", http.StatusForbidden)
		return
	}
	if !containsProtocol(r.Header.Get("Sec-WebSocket-Protocol"), eventSubprotocol) {
		http.Error(w, "required subprotocol missing", http.StatusBadRequest)
		return
	}
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{Subprotocols: []string{eventSubprotocol}, InsecureSkipVerify: true})
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	var request subscribeMessage
	err = wsjson.Read(ctx, c, &request)
	cancel()
	if err != nil || request.Type != "subscribe" || request.ProtocolVersion != 1 {
		_ = wsjson.Write(r.Context(), c, map[string]any{"type": "error", "code": "invalid_subscription", "message": "Send a version 1 subscribe message"})
		return
	}

	live, unsubscribe, status, replay, snapshot := s.Reducer.Bootstrap(request.LastSequence)
	defer unsubscribe()
	through := snapshot.Sequence
	cursor := request.LastSequence
	if request.LastSequence == 0 || status == "resync_required" {
		if status == "resync_required" {
			if wsjson.Write(r.Context(), c, map[string]any{"type": "resync_required", "cursor": through, "reason": "requested_sequence_outside_retention"}) != nil {
				return
			}
		}
		if wsjson.Write(r.Context(), c, map[string]any{"type": "snapshot", "cursor": through, "snapshot": snapshot}) != nil {
			return
		}
		cursor = through
	} else {
		for _, event := range replay {
			if wsjson.Write(r.Context(), c, map[string]any{"type": "event", "event": event}) != nil {
				return
			}
			cursor = event.Sequence
		}
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-live:
			if !ok {
				_ = c.Close(websocket.StatusPolicyViolation, "slow client")
				return
			}
			if event.Sequence <= cursor {
				continue
			}
			if wsjson.Write(r.Context(), c, map[string]any{"type": "event", "event": event}) != nil {
				return
			}
			cursor = event.Sequence
		case <-ticker.C:
			if wsjson.Write(r.Context(), c, map[string]any{"type": "heartbeat", "cursor": cursor}) != nil {
				return
			}
		}
	}
}

func containsProtocol(value, expected string) bool {
	for _, item := range strings.Split(value, ",") {
		if strings.TrimSpace(item) == expected {
			return true
		}
	}
	return false
}
