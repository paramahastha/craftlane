package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/craftlane/craftlane/internal/protocol"
)

const (
	DefaultRetention = 10000
	MinimumRetention = 1000
	segmentSize      = 1000
)

func (s Store) eventSegments() ([]string, error) {
	paths, err := filepath.Glob(s.path("events-*.jsonl"))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(s.path("events.jsonl")); err == nil {
		paths = append(paths, s.path("events.jsonl"))
	}
	sort.Slice(paths, func(i, j int) bool {
		if filepath.Base(paths[i]) == "events.jsonl" {
			return true
		}
		if filepath.Base(paths[j]) == "events.jsonl" {
			return false
		}
		return paths[i] < paths[j]
	})
	return paths, nil
}

func (s Store) AppendEvent(e protocol.Event) error {
	if err := os.MkdirAll(s.Dir, 0700); err != nil {
		return err
	}
	paths, err := s.eventSegments()
	if err != nil {
		return err
	}
	name := s.path("events.jsonl")
	if len(paths) > 0 {
		name = paths[len(paths)-1]
		count, countErr := countLines(name)
		if countErr != nil {
			return countErr
		}
		if count >= segmentSize {
			name = s.path(fmt.Sprintf("events-%06d.jsonl", len(paths)+1))
		}
	}
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = f.Write(append(b, '\n'))
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	return err
}

func countLines(name string) (int, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		n++
	}
	return n, sc.Err()
}

func (s Store) ReadEvents() ([]protocol.Event, error) {
	paths, err := s.eventSegments()
	if err != nil {
		return nil, err
	}
	var out []protocol.Event
	for _, name := range paths {
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			var e protocol.Event
			if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
				_ = f.Close()
				return nil, err
			}
			out = append(out, e)
		}
		err = sc.Err()
		_ = f.Close()
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// RemoveEvent rolls back the just-appended event if snapshot persistence fails.
func (s Store) RemoveEvent(sequence uint64) error {
	paths, err := s.eventSegments()
	if err != nil || len(paths) == 0 {
		return err
	}
	name := paths[len(paths)-1]
	b, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return nil
	}
	var last protocol.Event
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &last); err != nil {
		return err
	}
	if last.Sequence != sequence {
		return fmt.Errorf("event %d is not the latest event", sequence)
	}
	if len(lines) == 1 {
		return os.Remove(name)
	}
	return os.WriteFile(name, []byte(strings.Join(lines[:len(lines)-1], "\n")+"\n"), 0600)
}

func (s Store) ReplaceEvents(events []protocol.Event) error {
	// Compatibility repair for old stores; normal appends never rewrite history.
	if err := os.MkdirAll(s.Dir, 0700); err != nil {
		return err
	}
	paths, _ := s.eventSegments()
	for _, name := range paths {
		_ = os.Remove(name)
	}
	for _, event := range events {
		if err := s.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

// PruneEvents removes complete old segments and rewrites only the one segment
// crossing the retention boundary. It never rewrites the retained history.
func (s Store) PruneEvents(retention int) error {
	events, err := s.ReadEvents()
	if err != nil || len(events) <= retention {
		return err
	}
	cutoff := events[len(events)-retention].Sequence
	paths, err := s.eventSegments()
	if err != nil {
		return err
	}
	for _, name := range paths {
		segment, err := readSegment(name)
		if err != nil {
			return err
		}
		keep := segment[:0]
		for _, event := range segment {
			if event.Sequence >= cutoff {
				keep = append(keep, event)
			}
		}
		if len(keep) == 0 {
			if err := os.Remove(name); err != nil && !os.IsNotExist(err) {
				return err
			}
			continue
		}
		if len(keep) == len(segment) {
			continue
		}
		tmp := name + ".tmp"
		f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(f)
		for _, event := range keep {
			if err = enc.Encode(event); err != nil {
				_ = f.Close()
				return err
			}
		}
		if err = f.Sync(); err == nil {
			err = f.Close()
		} else {
			_ = f.Close()
		}
		if err != nil {
			return err
		}
		if err = os.Rename(tmp, name); err != nil {
			return err
		}
	}
	return nil
}

// TruncateEventsAfter discards journal entries that were written after the
// last committed snapshot, recovering a crash between append and snapshot.
func (s Store) TruncateEventsAfter(sequence uint64) error {
	paths, err := s.eventSegments()
	if err != nil {
		return err
	}
	for _, name := range paths {
		segment, err := readSegment(name)
		if err != nil {
			return err
		}
		keep := segment[:0]
		for _, event := range segment {
			if event.Sequence <= sequence {
				keep = append(keep, event)
			}
		}
		if len(keep) == len(segment) {
			continue
		}
		if len(keep) == 0 {
			if err := os.Remove(name); err != nil && !os.IsNotExist(err) {
				return err
			}
			continue
		}
		if err := writeSegment(name, keep); err != nil {
			return err
		}
	}
	return nil
}

func writeSegment(name string, events []protocol.Event) error {
	tmp := name + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	for _, event := range events {
		if err = enc.Encode(event); err != nil {
			_ = f.Close()
			return err
		}
	}
	if err = f.Sync(); err == nil {
		err = f.Close()
	} else {
		_ = f.Close()
	}
	if err != nil {
		return err
	}
	return os.Rename(tmp, name)
}

func readSegment(name string) ([]protocol.Event, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var events []protocol.Event
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var event protocol.Event
		if err := json.Unmarshal(sc.Bytes(), &event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, sc.Err()
}

func NormalizeRetention(retention int) int {
	if retention == 0 {
		return DefaultRetention
	}
	if retention < MinimumRetention {
		return MinimumRetention
	}
	return retention
}

func RetentionFromEnv(value string) int {
	if value == "" {
		return DefaultRetention
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return DefaultRetention
	}
	return NormalizeRetention(n)
}
