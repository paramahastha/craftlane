package storage

import (
	"encoding/json"
	"fmt"
	"github.com/craftlane/craftlane/internal/protocol"
	"os"
	"path/filepath"
)

type Store struct{ Dir string }

func New(dir string) Store              { return Store{Dir: dir} }
func (s Store) path(name string) string { return filepath.Join(s.Dir, name) }
func (s Store) write(name string, v any) error {
	if e := os.MkdirAll(s.Dir, 0700); e != nil {
		return e
	}
	b, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		return e
	}
	t := s.path(name + ".tmp")
	if e = os.WriteFile(t, b, 0600); e != nil {
		return e
	}
	return os.Rename(t, s.path(name))
}
func (s Store) SaveRecord(r protocol.RuntimeRecord) error { return s.write("runtime.json", r) }
func (s Store) LoadRecord() (protocol.RuntimeRecord, error) {
	var r protocol.RuntimeRecord
	b, e := os.ReadFile(s.path("runtime.json"))
	if e != nil {
		return r, e
	}
	e = json.Unmarshal(b, &r)
	return r, e
}
func (s Store) SaveSnapshot(v protocol.Snapshot) error { return s.write("snapshot.json", v) }
func (s Store) LoadSnapshot() (protocol.Snapshot, error) {
	var v protocol.Snapshot
	b, e := os.ReadFile(s.path("snapshot.json"))
	if e != nil {
		return v, e
	}
	e = json.Unmarshal(b, &v)
	return v, e
}
func (s Store) SaveControl(secret string) error {
	if e := os.MkdirAll(s.Dir, 0700); e != nil {
		return e
	}
	return os.WriteFile(s.path("control.token"), []byte(secret), 0600)
}
func (s Store) LoadControl() (string, error) {
	b, e := os.ReadFile(s.path("control.token"))
	return string(b), e
}
func (s Store) ClearControl() error {
	e := os.Remove(s.path("control.token"))
	if os.IsNotExist(e) {
		return nil
	}
	return e
}
func (s Store) Validate() error {
	if s.Dir == "" {
		return fmt.Errorf("empty storage directory")
	}
	return nil
}
