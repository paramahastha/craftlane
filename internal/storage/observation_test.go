package storage

import "testing"

func TestObservationSecretAndBoundedFingerprints(t *testing.T) {
	store := New(t.TempDir())
	if err := store.SaveObservationSecret("observation-secret"); err != nil {
		t.Fatal(err)
	}
	if got, err := store.LoadObservationSecret(); err != nil || got != "observation-secret" {
		t.Fatalf("secret = %q, %v", got, err)
	}
	for _, value := range []string{"a", "b", "c"} {
		if added, err := store.AddObservationFingerprint(value, 2); err != nil || !added {
			t.Fatalf("add %q = %v, %v", value, added, err)
		}
	}
	if added, err := store.AddObservationFingerprint("c", 2); err != nil || added {
		t.Fatalf("duplicate = %v, %v", added, err)
	}
	items, err := store.LoadObservationFingerprints()
	if err != nil || len(items) != 2 || items[0] != "b" || items[1] != "c" {
		t.Fatalf("fingerprints = %#v, %v", items, err)
	}
}
