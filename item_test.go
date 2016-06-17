package goque

import "testing"

func TestNewItemObjectAndUnmarshall(t *testing.T) {
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
	if err := item.Unmarshall(&b); err != nil {
		t.Error(err)
	}

	if a != b {
		t.Fail()
	}
}

func TestNewPriorityItemObjectAndUnmarshall(t *testing.T) {
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
		t.Fail()
	}

	var b object
	if err := item.Unmarshall(&b); err != nil {
		t.Error(err)
	}

	if a != b {
		t.Fail()
	}
}
