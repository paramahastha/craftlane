package app

import (
	"context"
	"fmt"
	"github.com/craftlane/craftlane/internal/diagnostics"
	"github.com/craftlane/craftlane/internal/events"
	"github.com/craftlane/craftlane/internal/project"
	"github.com/craftlane/craftlane/internal/protocol"
	"github.com/craftlane/craftlane/internal/server"
	"github.com/craftlane/craftlane/internal/storage"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Runtime struct {
	Project  protocol.Project
	Record   protocol.RuntimeRecord
	Store    storage.Store
	Token    string
	HTTP     *http.Server
	Listener net.Listener
	Reducer  *events.Reducer
	done     chan struct{}
	once     sync.Once
}

func stateDir(p protocol.Project) (string, error) {
	if base := os.Getenv("CRAFTLANE_STATE_DIR"); base != "" {
		return filepath.Join(base, p.ID), nil
	}
	base, e := os.UserConfigDir()
	if e != nil {
		return "", e
	}
	return filepath.Join(base, "craftlane", p.ID), nil
}
func StorageForProject(p protocol.Project) (storage.Store, error) {
	d, e := stateDir(p)
	if e != nil {
		return storage.Store{}, e
	}
	return storage.New(d), nil
}
func Start(ctx context.Context, p protocol.Project, port int) (*Runtime, error) {
	return StartWithID(ctx, p, port, "")
}

func StartWithID(ctx context.Context, p protocol.Project, port int, runtimeID string) (*Runtime, error) {
	dir, e := stateDir(p)
	if e != nil {
		return nil, e
	}
	st := storage.New(dir)
	token := storage.RandomSecret()
	ln, e := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
	if e != nil {
		return nil, fmt.Errorf("port %d unavailable: %w", port, e)
	}
	origin := fmt.Sprintf("http://%s", ln.Addr())
	id := runtimeID
	if id == "" {
		id = storage.RandomSecret()[:16]
	}
	rec := protocol.RuntimeRecord{RuntimeID: id, ProjectID: p.ID, ProjectRoot: p.Root, Origin: origin, PID: os.Getpid(), ProtocolVersion: 1, StartedAt: time.Now().UTC(), State: "starting"}
	if e = st.SaveRecord(rec); e != nil {
		ln.Close()
		return nil, e
	}
	rt := &Runtime{Project: p, Record: rec, Store: st, Token: token, Listener: ln, done: make(chan struct{})}
	ready := false
	defer func() {
		if ready {
			return
		}
		_ = ln.Close()
		_ = st.ClearControl()
		rec.State = "stopped"
		_ = st.SaveRecord(rec)
	}()
	if e = st.SaveControl(token); e != nil {
		return nil, e
	}
	pub := protocol.Runtime{ID: id, ProjectID: p.ID, Origin: origin, ProtocolVersion: 1, StartedAt: rec.StartedAt, State: "ready"}
	reducer := events.NewReducer(st, p, pub)
	previous, previousErr := st.LoadSnapshot()
	checks := diagnostics.Detect(p.Root, true, server.AssetReady())
	if e = reducer.SetDiagnostics(checks); e != nil {
		return nil, e
	}
	if previousErr == nil && project.Changed(previous.Project, p) {
		if _, e = reducer.Append("project.changed", map[string]any{"previous": previous.Project, "project": p}); e != nil {
			return nil, e
		}
	}
	if _, e = reducer.Append("project.detected", map[string]any{"project": p}); e != nil {
		return nil, e
	}
	if _, e = reducer.Append("runtime.ready", map[string]any{"runtimeId": id, "origin": origin}); e != nil {
		return nil, e
	}
	if e = ctx.Err(); e != nil {
		return nil, e
	}
	rec.State = "ready"
	if e = st.SaveRecord(rec); e != nil {
		return nil, e
	}
	rt.Reducer = reducer
	hs := server.New(p, pub, token)
	hs.Reducer = reducer
	hs.Checks = checks
	hs.Stop = func() { _ = rt.Shutdown(context.Background()) }
	rt.HTTP = &http.Server{Handler: hs.Handler()}
	go func() {
		if err := rt.HTTP.Serve(ln); err != nil && err != http.ErrServerClosed {
			_ = rt.Shutdown(context.Background())
		}
	}()
	go func() {
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = rt.Shutdown(shutdownCtx)
		case <-rt.done:
		}
	}()
	go monitorProjectChanges(ctx, rt)
	ready = true
	return rt, nil
}

func monitorProjectChanges(ctx context.Context, rt *Runtime) {
	check := func() (string, error) {
		checkCtx, cancel := context.WithTimeout(ctx, 750*time.Millisecond)
		defer cancel()
		return project.Revision(checkCtx, rt.Project.Root)
	}
	last, err := check()
	if err != nil {
		return
	}
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-rt.done:
			return
		case <-ticker.C:
			revision, err := check()
			if err != nil || revision == last {
				continue
			}
			last = revision
			_, _ = rt.Reducer.Append("project.changed", map[string]any{"project": rt.Project, "revision": revision})
		}
	}
}
func (r *Runtime) Shutdown(ctx context.Context) error {
	var result error
	r.once.Do(func() {
		r.Record.State = "stopping"
		_ = r.Store.SaveRecord(r.Record)
		if r.Reducer != nil {
			_ = r.Reducer.SetRuntimeState("stopping")
			_, _ = r.Reducer.Append("runtime.stopping", map[string]any{"runtimeId": r.Record.RuntimeID})
		}
		if r.HTTP != nil {
			result = r.HTTP.Shutdown(ctx)
		}
		if r.Listener != nil {
			_ = r.Listener.Close()
		}
		_ = r.Store.ClearControl()
		r.Record.State = "stopped"
		if err := r.Store.SaveRecord(r.Record); result == nil {
			result = err
		}
		close(r.done)
	})
	return result
}

func (r *Runtime) Done() <-chan struct{} { return r.done }
