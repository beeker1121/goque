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
