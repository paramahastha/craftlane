package project

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/craftlane/craftlane/internal/protocol"
	"path/filepath"
	"strings"
	"time"
)

func Identity(root, common string) string {
	h := sha256.Sum256([]byte(root + "\x00" + common))
	return hex.EncodeToString(h[:])
}
func New(root, common string) protocol.Project {
	return protocol.Project{ID: Identity(root, common), Name: filepath.Base(root), Root: root, GitCommonDir: common, DetectedAt: time.Now().UTC(), Capabilities: protocol.Capabilities{"git": "available", "speckit": "unavailable", "provider": "unavailable"}}
}

func Changed(previous, current protocol.Project) bool {
	return previous.ID != current.ID || previous.Name != current.Name || previous.Root != current.Root || previous.GitCommonDir != current.GitCommonDir
}
func IsWithin(root, path string) bool {
	r, e := filepath.Abs(root)
	if e != nil {
		return false
	}
	p, e := filepath.Abs(path)
	if e != nil {
		return false
	}
	rel, e := filepath.Rel(r, p)
	return e == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
