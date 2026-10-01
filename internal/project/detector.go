package project

import (
	"context"
	"fmt"
	"github.com/craftlane/craftlane/internal/protocol"
	"os/exec"
	"path/filepath"
	"strings"
)

func run(ctx context.Context, dir string, args ...string) (string, error) {
	c := exec.CommandContext(ctx, "git", args...)
	c.Dir = dir
	b, e := c.Output()
	if e != nil {
		return "", e
	}
	return strings.TrimSpace(string(b)), nil
}

// Revision returns the current Git revision used by the lightweight runtime
// change detector. It is intentionally bounded by the caller's context.
func Revision(ctx context.Context, dir string) (string, error) {
	return run(ctx, dir, "rev-parse", "HEAD")
}
func Detect(ctx context.Context, dir string) (protocol.Project, error) {
	d, e := filepath.Abs(dir)
	if e != nil {
		return protocol.Project{}, e
	}
	root, e := run(ctx, d, "rev-parse", "--show-toplevel")
	if e != nil {
		return protocol.Project{}, fmt.Errorf("not a Git repository: %w", e)
	}
	root, e = filepath.EvalSymlinks(root)
	if e != nil {
		return protocol.Project{}, e
	}
	common, e := run(ctx, root, "rev-parse", "--git-common-dir")
	if e != nil {
		return protocol.Project{}, e
	}
	if !filepath.IsAbs(common) {
		common = filepath.Join(root, common)
	}
	common, _ = filepath.EvalSymlinks(common)
	return New(root, common), nil
}
