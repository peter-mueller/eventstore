package eventstore

import "testing"

func TestHasIndex(t *testing.T) {
	var normal Information
	if !normal.HasIndex(0) {
		t.Errorf("normal must have index 0")
	}
	if !normal.HasIndex(1) {
		t.Errorf("normal must have index 1")
	}
	if !normal.HasIndex(10000) {
		t.Errorf("normal must have index 10000")
	}

	snapshot := Information{
		StartIndex: 100,
	}
	if snapshot.HasIndex(0) {
		t.Errorf("snapshot must not have index 0")
	}
	if snapshot.HasIndex(99) {
		t.Errorf("snapshot must not have index 99")
	}
	if !snapshot.HasIndex(100) {
		t.Errorf("snapshot must have index 100")
	}
	if !snapshot.HasIndex(101) {
		t.Errorf("snapshot must have index 101")
	}
	if !snapshot.HasIndex(100000) {
		t.Errorf("snapshot must have index 10000")
	}
}