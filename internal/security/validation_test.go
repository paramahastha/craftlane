package security

import "testing"

func TestValidation(t *testing.T) {
	if !Host("127.0.0.1:7331", "127.0.0.1:7331") || Host("localhost:7331", "127.0.0.1:7331") {
		t.Fatal("host")
	}
	if !Token("a", "a") || Token("a", "b") {
		t.Fatal("token")
	}
	if !ApprovedPath("/tmp/project", "/tmp/project/file") || ApprovedPath("/tmp/project", "/tmp/project/../secret") {
		t.Fatal("path")
	}
}
