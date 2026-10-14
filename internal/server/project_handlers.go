package server

import (
	"encoding/json"
	"github.com/craftlane/craftlane/internal/events"
	"github.com/craftlane/craftlane/internal/protocol"
	"net/http"
)

func jsonWrite(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "response encoding failed", http.StatusInternalServerError)
		return
	}
	var plain any
	if json.Unmarshal(b, &plain) != nil {
		http.Error(w, "response encoding failed", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(events.Redact(plain))
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	status := "ok"
	for _, check := range s.Checks {
		if check.Status == "fail" {
			status = "degraded"
			break
		}
	}
	snapshot := protocol.Snapshot{}
	if s.Reducer != nil {
		snapshot = s.Reducer.Snapshot()
	}
	jsonWrite(w, protocol.Health{Status: status, Runtime: s.Runtime, Project: s.Project, Checks: s.Checks, Agents: snapshot.Agents, ProviderStatuses: snapshot.ProviderStatuses})
}

func (s *Server) project(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	jsonWrite(w, s.Project)
}
