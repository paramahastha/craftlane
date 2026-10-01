package server

import (
	"net/http"
	"strconv"
)

func (s *Server) eventsHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if s.Reducer == nil {
		jsonWrite(w, map[string]any{"status": "ok", "cursor": uint64(0), "events": []any{}})
		return
	}
	after, err := strconv.ParseUint(r.URL.Query().Get("after"), 10, 64)
	if r.URL.Query().Get("after") == "" {
		after, err = 0, nil
	}
	if err != nil {
		http.Error(w, "invalid after cursor", http.StatusBadRequest)
		return
	}
	status, retained, snapshot := s.Reducer.ReplayCurrent(after)
	response := map[string]any{"status": status, "cursor": snapshot.Sequence, "events": retained}
	if status == "resync_required" {
		response["snapshot"] = snapshot
	}
	jsonWrite(w, response)
}
