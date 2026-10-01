package storage

import (
	"github.com/craftlane/craftlane/internal/protocol"
	"testing"
)

func TestAtomicRecordAndOwnerControl(t *testing.T) {
	s := New(t.TempDir())
	r := protocol.RuntimeRecord{RuntimeID: "r", State: "ready"}
	if err := s.SaveRecord(r); err != nil {
		t.Fatal(err)
	}
	got, err := s.LoadRecord()
	if err != nil || got.RuntimeID != "r" {
		t.Fatal(got, err)
	}
	if err = s.SaveControl("secret"); err != nil {
		t.Fatal(err)
	}
	if v, _ := s.LoadControl(); v != "secret" {
		t.Fatal(v)
	}
}
