package events

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/storage"
	"os"
	"sync"
	"time"
)

type Reducer struct {
	mu          sync.Mutex
	Store       storage.Store
	Project     protocol.Project
	Runtime     protocol.Runtime
	Events      []protocol.Event
	Sequence    uint64
	Retention   int
	Hub         *Hub
	diagnostics []protocol.DiagnosticCheck
}

func NewReducer(s storage.Store, p protocol.Project, r protocol.Runtime) *Reducer {
	return NewReducerWithRetention(s, p, r, storage.RetentionFromEnv(os.Getenv("CRAFTLANE_EVENT_RETENTION")))
}

func NewReducerWithRetention(s storage.Store, p protocol.Project, runtime protocol.Runtime, retention int) *Reducer {
	r := &Reducer{Store: s, Project: p, Runtime: runtime, Retention: storage.NormalizeRetention(retention), Hub: NewHub(64)}
	snapshot, snapshotErr := s.LoadSnapshot()
	if events, err := s.ReadEvents(); err == nil {
		if snapshotErr == nil && len(events) > 0 && snapshot.Sequence < events[len(events)-1].Sequence {
			_ = s.TruncateEventsAfter(snapshot.Sequence)
			events, _ = s.ReadEvents()
		}
		r.Events = events
		if len(events) > r.Retention {
			r.Events = append([]protocol.Event(nil), events[len(events)-r.Retention:]...)
			_ = s.PruneEvents(r.Retention)
		}
		if len(r.Events) > 0 {
			r.Sequence = r.Events[len(r.Events)-1].Sequence
		}
	}
	if snapshotErr == nil && snapshot.Sequence > r.Sequence {
		r.Sequence = snapshot.Sequence
	}
	return r
}

func id() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *Reducer) Append(typ string, payload map[string]any) (protocol.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	eventID, err := id()
	if err != nil {
		return protocol.Event{}, err
	}
	redacted, ok := Redact(payload).(map[string]any)
	if !ok {
		return protocol.Event{}, fmt.Errorf("event payload must be an object")
	}
	previous := append([]protocol.Event(nil), r.Events...)
	r.Sequence++
	e := protocol.Event{SchemaVersion: 1, EventID: eventID, Sequence: r.Sequence, Timestamp: time.Now().UTC(), ProjectID: r.Project.ID, Source: protocol.EventSource{Name: "core"}, Type: typ, Payload: redacted}
	r.Events = append(r.Events, e)
	if err := r.Store.AppendEvent(e); err != nil {
		r.Sequence--
		r.Events = previous
		return protocol.Event{}, err
	}
	if err := r.Store.SaveSnapshot(r.snapshotLocked()); err != nil {
		_ = r.Store.RemoveEvent(e.Sequence)
		r.Sequence--
		r.Events = previous
		return protocol.Event{}, err
	}
	if len(r.Events) > r.Retention {
		r.Events = r.Events[len(r.Events)-r.Retention:]
		_ = r.Store.PruneEvents(r.Retention)
	}
	r.Hub.Publish(e)
	return e, nil
}
func (r *Reducer) Snapshot() protocol.Snapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.snapshotLocked()
}

func (r *Reducer) snapshotLocked() protocol.Snapshot {
	return protocol.Snapshot{Sequence: r.Sequence, Project: r.Project, Runtime: r.Runtime, Diagnostics: append([]protocol.DiagnosticCheck(nil), r.diagnostics...)}
}

func (r *Reducer) SetDiagnostics(checks []protocol.DiagnosticCheck) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.diagnostics = redactChecks(checks)
	return r.Store.SaveSnapshot(r.snapshotLocked())
}

func redactChecks(checks []protocol.DiagnosticCheck) []protocol.DiagnosticCheck {
	out := make([]protocol.DiagnosticCheck, len(checks))
	for i, check := range checks {
		check.Detail = RedactText(check.Detail)
		check.Action = RedactText(check.Action)
		out[i] = check
	}
	return out
}

func (r *Reducer) Replay(after, through uint64) (string, []protocol.Event, protocol.Snapshot) {
	r.mu.Lock()
	defer r.mu.Unlock()
	snapshot := r.snapshotLocked()
	if after > r.Sequence || (len(r.Events) > 0 && after+1 < r.Events[0].Sequence) {
		return "resync_required", nil, snapshot
	}
	out := make([]protocol.Event, 0)
	for _, event := range r.Events {
		if event.Sequence > after && event.Sequence <= through {
			out = append(out, event)
		}
	}
	return "ok", out, snapshot
}

// Bootstrap registers before taking the replay boundary, so appends cannot be
// missed between subscription and replay.
func (r *Reducer) Bootstrap(after uint64) (<-chan protocol.Event, func(), string, []protocol.Event, protocol.Snapshot) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, cancel := r.Hub.Subscribe()
	status, replay, snapshot := r.replayLocked(after, r.Sequence)
	return ch, cancel, status, replay, snapshot
}

func (r *Reducer) ReplayCurrent(after uint64) (string, []protocol.Event, protocol.Snapshot) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.replayLocked(after, r.Sequence)
}

func (r *Reducer) replayLocked(after, through uint64) (string, []protocol.Event, protocol.Snapshot) {
	snapshot := r.snapshotLocked()
	if after > through || (len(r.Events) > 0 && after+1 < r.Events[0].Sequence) {
		return "resync_required", nil, snapshot
	}
	out := make([]protocol.Event, 0)
	for _, event := range r.Events {
		if event.Sequence > after && event.Sequence <= through {
			out = append(out, event)
		}
	}
	return "ok", out, snapshot
}

func (r *Reducer) Subscribe() (<-chan protocol.Event, func(), uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, cancel := r.Hub.Subscribe()
	return ch, cancel, r.Sequence
}

func (r *Reducer) SetRuntimeState(state string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Runtime.State = state
	return r.Store.SaveSnapshot(r.snapshotLocked())
}
