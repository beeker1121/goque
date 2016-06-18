package goque

import "testing"

func TestNewItemObject(t *testing.T) {
	type object struct {
		X int
		Y int
	}

	a := object{X: 1, Y: 2}
	item, err := NewItemObject(a)
	if err != nil {
		t.Error(err)
	}

	var b object
	if err := item.ToObject(&b); err != nil {
		t.Error(err)
	}

	if a != b {
		t.Error("Decoded object does not match original object")
	}
}

func TestNewPriorityItemObject(t *testing.T) {
	type object struct {
		X int
		Y int
	}

	a := object{X: 1, Y: 2}
	item, err := NewPriorityItemObject(a, 42)
	if err != nil {
		t.Error(err)
	}

	if item.Priority != 42 {
		t.Errorf("Expected priority level to be 42, got %d", item.Priority)
	}

	var b object
	if err := item.ToObject(&b); err != nil {
		t.Error(err)
	}

	if a != b {
		t.Error("Decoded object does not match original object")
	}
}
