package events

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/storage"
	"os"
	"sort"
	"sync"
	"time"
)

type Reducer struct {
	mu               sync.Mutex
	Store            storage.Store
	Project          protocol.Project
	Runtime          protocol.Runtime
	Events           []protocol.Event
	Sequence         uint64
	Retention        int
	Hub              *Hub
	diagnostics      []protocol.DiagnosticCheck
	agents           map[string]protocol.Agent
	providerStatuses []protocol.ProviderStatus
	fingerprints     []string
}

func NewReducer(s storage.Store, p protocol.Project, r protocol.Runtime) *Reducer {
	return NewReducerWithRetention(s, p, r, storage.RetentionFromEnv(os.Getenv("CRAFTLANE_EVENT_RETENTION")))
}

func NewReducerWithRetention(s storage.Store, p protocol.Project, runtime protocol.Runtime, retention int) *Reducer {
	r := &Reducer{Store: s, Project: p, Runtime: runtime, Retention: storage.NormalizeRetention(retention), Hub: NewHub(64), agents: make(map[string]protocol.Agent)}
	snapshot, snapshotErr := s.LoadSnapshot()
	if snapshotErr == nil {
		for _, agent := range snapshot.Agents {
			agent.State = protocol.AgentUnknown
			agent.StateSource = "reported"
			agent.StaleAfterRestart = true
			r.agents[agent.ID] = agent
		}
		r.providerStatuses = append([]protocol.ProviderStatus(nil), snapshot.ProviderStatuses...)
		r.diagnostics = append([]protocol.DiagnosticCheck(nil), snapshot.Diagnostics...)
	}
	if fingerprints, err := s.LoadObservationFingerprints(); err == nil {
		r.fingerprints = fingerprints
	}
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
	return r.appendLocked(typ, payload, protocol.EventSource{Name: "core"})
}

func (r *Reducer) appendLocked(typ string, payload map[string]any, source protocol.EventSource) (protocol.Event, error) {
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
	e := protocol.Event{SchemaVersion: 1, EventID: eventID, Sequence: r.Sequence, Timestamp: time.Now().UTC(), ProjectID: r.Project.ID, Source: source, Type: typ, Payload: redacted}
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
	agents := make([]protocol.Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}
	sort.Slice(agents, func(i, j int) bool { return agents[i].ID < agents[j].ID })
	return protocol.Snapshot{Sequence: r.Sequence, Project: r.Project, Runtime: r.Runtime, Diagnostics: append([]protocol.DiagnosticCheck(nil), r.diagnostics...), Agents: agents, ProviderStatuses: append([]protocol.ProviderStatus(nil), r.providerStatuses...)}
}

func (r *Reducer) SetProviderStatuses(statuses []protocol.ProviderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providerStatuses = append([]protocol.ProviderStatus(nil), statuses...)
	return r.Store.SaveSnapshot(r.snapshotLocked())
}

func (r *Reducer) Agents() []protocol.Agent {
	r.mu.Lock()
	defer r.mu.Unlock()
	snapshot := r.snapshotLocked()
	return snapshot.Agents
}

func (r *Reducer) ApplyObservation(observation protocol.ObservationEnvelope, source protocol.EventSource) (protocol.Agent, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if observation.ProjectID != r.Project.ID {
		return protocol.Agent{}, false, fmt.Errorf("observation project does not match runtime")
	}
	if observation.Provider == "" || observation.Kind == "" || (observation.AgentKey == "" && observation.SessionKey == "") {
		return protocol.Agent{}, false, fmt.Errorf("invalid observation")
	}
	fingerprint := ObservationFingerprint(observation)
	for _, item := range r.fingerprints {
		if item == fingerprint {
			return protocol.Agent{}, false, nil
		}
	}
	id := ObservationAgentID(observation)
	previous, existed := r.agents[id]
	oldSnapshot := r.snapshotLocked()
	agent, changed, err := applyObservation(previous, observation)
	if err != nil {
		return protocol.Agent{}, false, err
	}
	r.fingerprints = append(r.fingerprints, fingerprint)
	if len(r.fingerprints) > storage.MinimumRetention {
		r.fingerprints = r.fingerprints[len(r.fingerprints)-storage.MinimumRetention:]
	}
	if !changed {
		if err := r.Store.SaveObservationFingerprints(r.fingerprints); err != nil {
			return protocol.Agent{}, false, err
		}
		return previous, false, nil
	}
	r.agents[id] = agent
	oldSequence, oldEvents, oldFingerprints := r.Sequence, r.Events, append([]string(nil), r.fingerprints...)
	payload := map[string]any{"agent": agent}
	if _, err := r.appendLocked("agent.updated", payload, source); err != nil {
		r.Sequence, r.Events, r.fingerprints = oldSequence, oldEvents, oldFingerprints
		if existed {
			r.agents[id] = previous
		} else {
			delete(r.agents, id)
		}
		return protocol.Agent{}, false, err
	}
	if err := r.Store.SaveObservationFingerprints(r.fingerprints); err != nil {
		_ = r.Store.RemoveEvent(r.Sequence)
		r.Sequence, r.Events, r.fingerprints = oldSequence, oldEvents, oldFingerprints
		if existed {
			r.agents[id] = previous
		} else {
			delete(r.agents, id)
		}
		_ = r.Store.SaveSnapshot(oldSnapshot)
		return protocol.Agent{}, false, err
	}
	return agent, true, nil
}

func (r *Reducer) InferIdle(now time.Time, threshold time.Duration) ([]protocol.Agent, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if threshold <= 0 {
		return nil, nil
	}
	var changed []protocol.Agent
	for id, agent := range r.agents {
		if (agent.State != protocol.AgentWorking && agent.State != protocol.AgentUnknown) || agent.LastObservedAt.IsZero() || now.Sub(agent.LastObservedAt) < threshold {
			continue
		}
		agent.State, agent.StateSource = protocol.AgentIdle, "inferred"
		r.agents[id] = agent
		if _, err := r.appendLocked("agent.updated", map[string]any{"agent": agent}, protocol.EventSource{Name: agent.Provider}); err != nil {
			return changed, err
		}
		changed = append(changed, agent)
	}
	return changed, nil
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
