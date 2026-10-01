package events

import (
	"github.com/craftlane/craftlane/internal/protocol"
	"sync"
)

type Hub struct {
	mu       sync.Mutex
	clients  map[chan protocol.Event]struct{}
	capacity int
}

func NewHub(capacity int) *Hub {
	if capacity < 1 {
		capacity = 1
	}
	return &Hub{clients: make(map[chan protocol.Event]struct{}), capacity: capacity}
}

func (h *Hub) Subscribe() (<-chan protocol.Event, func()) {
	h.mu.Lock()
	ch := make(chan protocol.Event, h.capacity)
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch, func() {
		h.mu.Lock()
		if _, ok := h.clients[ch]; ok {
			delete(h.clients, ch)
			close(ch)
		}
		h.mu.Unlock()
	}
}

func (h *Hub) Publish(event protocol.Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- event:
		default:
			delete(h.clients, ch)
			close(ch)
		}
	}
}
