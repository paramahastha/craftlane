package project

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDetectGitRoot(t *testing.T) {
	d := t.TempDir()
	if err := exec.Command("git", "init", d).Run(); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(d, "nested")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	p, err := Detect(context.Background(), sub)
	if err != nil {
		t.Fatal(err)
	}
	canonical, _ := filepath.EvalSymlinks(d)
	if p.Root != canonical || p.ID == "" {
		t.Fatalf("bad project: %+v", p)
	}
}
func TestDetectRejectsNonGit(t *testing.T) {
	_, err := Detect(context.Background(), t.TempDir())
	if err == nil {
		t.Fatal("expected rejection")
	}
}

func TestRevisionChangesAfterCommit(t *testing.T) {
	dir := t.TempDir()
	if _, err := run(context.Background(), dir, "init"); err != nil {
		t.Fatal(err)
	}
	if _, err := run(context.Background(), dir, "config", "user.email", "test@example.com"); err != nil {
		t.Fatal(err)
	}
	if _, err := run(context.Background(), dir, "config", "user.name", "Test"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "file"), []byte("one"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := run(context.Background(), dir, "add", "."); err != nil {
		t.Fatal(err)
	}
	if _, err := run(context.Background(), dir, "commit", "-m", "one"); err != nil {
		t.Fatal(err)
	}
	first, err := Revision(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "file"), []byte("two"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := run(context.Background(), dir, "add", "."); err != nil {
		t.Fatal(err)
	}
	if _, err := run(context.Background(), dir, "commit", "-m", "two"); err != nil {
		t.Fatal(err)
	}
	second, err := Revision(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("revision did not change")
	}
}
