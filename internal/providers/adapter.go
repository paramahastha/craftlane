package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func executableStatus(name string) (string, string) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "unavailable", name + " executable was not found"
	}
	version := "available"
	if out, err := exec.Command(path, "--version").Output(); err == nil && len(out) > 0 {
		version = strings.TrimSpace(string(out))
		if len(version) > 80 {
			version = version[:80]
		}
	}
	return "ready", version
}

func ExecutableStatusForTests(name string) (string, string) { return executableStatus(name) }

func readObject(path string) (map[string]any, error) {
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]any{}, nil
	}
	if err != nil {
		return nil, err
	}
	var object map[string]any
	if err := json.Unmarshal(b, &object); err != nil {
		return nil, fmt.Errorf("invalid provider configuration: %w", err)
	}
	return object, nil
}

func writeObject(path string, value map[string]any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func deleteIfEmpty(path string, object map[string]any) error {
	if len(object) == 0 {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
