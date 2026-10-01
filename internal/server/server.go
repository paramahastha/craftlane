package server

import (
	"github.com/craftlane/craftlane/internal/events"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/security"
	"io/fs"
	"net/http"
	"strings"
)

type Server struct {
	Project        protocol.Project
	Runtime        protocol.Runtime
	Reducer        *events.Reducer
	Checks         []protocol.DiagnosticCheck
	ExpectedHost   string
	ExpectedOrigin string
	Token          string
	Stop           func()
}

func New(p protocol.Project, r protocol.Runtime, token string) *Server {
	return &Server{Project: p, Runtime: r, Token: token, ExpectedHost: strings.TrimPrefix(r.Origin, "http://"), ExpectedOrigin: r.Origin}
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", s.health)
	mux.HandleFunc("/api/v1/project", s.project)
	mux.HandleFunc("/api/v1/events", s.eventsHTTP)
	mux.HandleFunc("/ws", s.websocket)
	mux.HandleFunc("/api/v1/runtime/stop", s.stop)
	sub, _ := fs.Sub(assetFS, "assets/dist")
	static := http.FileServer(http.FS(sub))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && !strings.HasPrefix(r.URL.Path, "/assets/") {
			r.URL.Path = "/"
		}
		static.ServeHTTP(w, r)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !security.Host(r.Host, s.ExpectedHost) {
			http.Error(w, "invalid host", http.StatusBadRequest)
			return
		}
		mux.ServeHTTP(w, r)
	})
}
func (s *Server) stop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Origin") != "" && !security.Origin(r.Header.Get("Origin"), s.ExpectedOrigin) {
		http.Error(w, "forbidden origin", http.StatusForbidden)
		return
	}
	if !security.Token(r.Header.Get("X-Craftlane-Token"), s.Token) {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if s.Stop != nil {
		go s.Stop()
	}
}
