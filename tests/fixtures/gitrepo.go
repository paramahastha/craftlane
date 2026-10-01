package fixtures

import (
	"os/exec"
	"testing"
)

func GitRepo(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	if err := exec.Command("git", "init", d).Run(); err != nil {
		t.Fatal(err)
	}
	return d
}
