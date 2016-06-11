package goque

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestStackDrop(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Error(err)
	}

	s.Drop()

	if _, err := os.Stat(file); err == nil {
		t.Error("Expected directory for test database to have been deleted")
	}
}

func TestStackPush(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := s.Push(item); err != nil {
			t.Error(err)
		}
	}

	if s.Length() != 10 {
		t.Errorf("Expected queue size of 10, got %d", s.Length())
	}
}

func TestStackPop(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := s.Push(item); err != nil {
			t.Error(err)
		}
	}

	if s.Length() != 10 {
		t.Errorf("Expected queue length of 1, got %d", s.Length())
	}

	deqItem, err := s.Pop()
	if err != nil {
		t.Error(err)
	}

	if s.Length() != 9 {
		t.Errorf("Expected queue length of 0, got %d", s.Length())
	}

	compStr := "value for item 10"

	if deqItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, deqItem.ToString())
	}
}

func TestStackPeek(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	compStr := "value for item"

	if err := s.Push(NewItemString(compStr)); err != nil {
		t.Error(err)
	}

	peekItem, err := s.Peek()
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}
}

func TestStackPeekByOffset(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := s.Push(item); err != nil {
			t.Error(err)
		}
	}

	peekItem, err := s.PeekByOffset(3)
	if err != nil {
		t.Error(err)
	}

	compStr := "value for item 7"

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}
}

func TestStackPeekByID(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := s.Push(item); err != nil {
			t.Error(err)
		}
	}

	peekItem, err := s.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	compStr := "value for item 3"

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}
}

func TestStackUpdate(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := s.Push(item); err != nil {
			t.Error(err)
		}
	}

	item, err := s.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	oldCompStr := "value for item 3"
	newCompStr := "new value for item 3"

	if item.ToString() != oldCompStr {
		t.Errorf("Expected string to be '%s', got '%s'", oldCompStr, item.ToString())
	}

	if err := s.Update(item, []byte(newCompStr)); err != nil {
		t.Error(err)
	}

	if item.ToString() != newCompStr {
		t.Errorf("Expected current item value to be '%s', got '%s'", newCompStr, item.ToString())
	}

	newItem, err := s.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	if newItem.ToString() != newCompStr {
		t.Errorf("Expected new item value to be '%s', got '%s'", newCompStr, item.ToString())
	}
}

func TestStackUpdateString(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := s.Push(item); err != nil {
			t.Error(err)
		}
	}

	item, err := s.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	oldCompStr := "value for item 3"
	newCompStr := "new value for item 3"

	if item.ToString() != oldCompStr {
		t.Errorf("Expected string to be '%s', got '%s'", oldCompStr, item.ToString())
	}

	if err := s.UpdateString(item, newCompStr); err != nil {
		t.Error(err)
	}

	if item.ToString() != newCompStr {
		t.Errorf("Expected current item value to be '%s', got '%s'", newCompStr, item.ToString())
	}

	newItem, err := s.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	if newItem.ToString() != newCompStr {
		t.Errorf("Expected new item value to be '%s', got '%s'", newCompStr, item.ToString())
	}
}

func TestStackEmpty(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	err = s.Push(NewItemString("value for item"))
	if err != nil {
		t.Error(err)
	}

	_, err = s.Pop()
	if err != nil {
		t.Error(err)
	}

	_, err = s.Pop()
	if err != ErrEmpty {
		t.Errorf("Expected to get queue empty error, got %s", err.Error())
	}
}

func TestStackOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	err = s.Push(NewItemString("value for item"))
	if err != nil {
		t.Error(err)
	}

	_, err = s.PeekByOffset(2)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}
}

func BenchmarkStackEnqueue(b *testing.B) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		b.Error(err)
	}
	defer s.Drop()

	for n := 0; n < b.N; n++ {
		if err := s.Push(NewItemString("value")); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkStackDequeue(b *testing.B) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		b.Error(err)
	}
	defer s.Drop()

	for n := 0; n < b.N; n++ {
		for i := 1; i <= 10; i++ {
			if err := s.Push(NewItemString("value")); err != nil {
				b.Error(err)
			}
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = s.Pop()
	}
}
