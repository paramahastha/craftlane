package security

import (
	"crypto/subtle"
	"net"
	"net/url"
	"path/filepath"
	"strings"
)

func Host(host, expected string) bool     { return host == expected }
func Origin(origin, expected string) bool { return origin == expected }
func Token(got, expected string) bool {
	return got != "" && expected != "" && len(got) == len(expected) && subtle.ConstantTimeCompare([]byte(got), []byte(expected)) == 1
}
func ApprovedPath(root, path string) bool {
	r, e := canonicalPath(root)
	if e != nil {
		return false
	}
	p, e := canonicalPath(path)
	if e != nil {
		return false
	}
	rel, e := filepath.Rel(r, p)
	return e == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func canonicalPath(path string) (string, error) {
	current, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	var suffix []string
	for {
		resolved, evalErr := filepath.EvalSymlinks(current)
		if evalErr == nil {
			return filepath.Join(append([]string{resolved}, suffix...)...), nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", evalErr
		}
		suffix = append([]string{filepath.Base(current)}, suffix...)
		current = parent
	}
}
func Loopback(addr string) bool {
	h, _, e := net.SplitHostPort(addr)
	return e == nil && h == "127.0.0.1"
}
func ExactOrigin(value string) bool {
	u, e := url.Parse(value)
	return e == nil && u.Scheme == "http" && u.Hostname() == "127.0.0.1" && u.Path == ""
}
