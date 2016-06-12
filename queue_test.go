package goque

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestQueueDrop(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}

	if _, err = os.Stat(file); os.IsNotExist(err) {
		t.Error(err)
	}

	q.Drop()

	if _, err = os.Stat(file); err == nil {
		t.Error("Expected directory for test database to have been deleted")
	}
}

func TestQueueEnqueue(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err = q.Enqueue(item); err != nil {
			t.Error(err)
		}
	}

	if q.Length() != 10 {
		t.Errorf("Expected queue size of 10, got %d", q.Length())
	}
}

func TestQueueDequeue(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := q.Enqueue(item); err != nil {
			t.Error(err)
		}
	}

	if q.Length() != 10 {
		t.Errorf("Expected queue length of 1, got %d", q.Length())
	}

	deqItem, err := q.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if q.Length() != 9 {
		t.Errorf("Expected queue length of 0, got %d", q.Length())
	}

	compStr := "value for item 1"

	if deqItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, deqItem.ToString())
	}
}

func TestQueuePeek(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	compStr := "value for item"

	if err := q.Enqueue(NewItemString(compStr)); err != nil {
		t.Error(err)
	}

	peekItem, err := q.Peek()
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}
}

func TestQueuePeekByOffset(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := q.Enqueue(item); err != nil {
			t.Error(err)
		}
	}

	peekItem, err := q.PeekByOffset(3)
	if err != nil {
		t.Error(err)
	}

	compStr := "value for item 3"

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}
}

func TestQueuePeekByID(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := q.Enqueue(item); err != nil {
			t.Error(err)
		}
	}

	peekItem, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	compStr := "value for item 3"

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}
}

func TestQueueUpdate(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := q.Enqueue(item); err != nil {
			t.Error(err)
		}
	}

	item, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	oldCompStr := "value for item 3"
	newCompStr := "new value for item 3"

	if item.ToString() != oldCompStr {
		t.Errorf("Expected string to be '%s', got '%s'", oldCompStr, item.ToString())
	}

	if err := q.Update(item, []byte(newCompStr)); err != nil {
		t.Error(err)
	}

	if item.ToString() != newCompStr {
		t.Errorf("Expected current item value to be '%s', got '%s'", newCompStr, item.ToString())
	}

	newItem, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	if newItem.ToString() != newCompStr {
		t.Errorf("Expected new item value to be '%s', got '%s'", newCompStr, item.ToString())
	}
}

func TestQueueUpdateString(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	for i := 1; i <= 10; i++ {
		item := NewItemString(fmt.Sprintf("value for item %d", i))
		if err := q.Enqueue(item); err != nil {
			t.Error(err)
		}
	}

	item, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	oldCompStr := "value for item 3"
	newCompStr := "new value for item 3"

	if item.ToString() != oldCompStr {
		t.Errorf("Expected string to be '%s', got '%s'", oldCompStr, item.ToString())
	}

	if err := q.UpdateString(item, newCompStr); err != nil {
		t.Error(err)
	}

	if item.ToString() != newCompStr {
		t.Errorf("Expected current item value to be '%s', got '%s'", newCompStr, item.ToString())
	}

	newItem, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	if newItem.ToString() != newCompStr {
		t.Errorf("Expected new item value to be '%s', got '%s'", newCompStr, item.ToString())
	}
}

func TestQueueEmpty(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	err = q.Enqueue(NewItemString("value for item"))
	if err != nil {
		t.Error(err)
	}

	_, err = q.Dequeue()
	if err != nil {
		t.Error(err)
	}

	_, err = q.Dequeue()
	if err != ErrEmpty {
		t.Errorf("Expected to get queue empty error, got %s", err.Error())
	}
}

func TestQueueOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	err = q.Enqueue(NewItemString("value for item"))
	if err != nil {
		t.Error(err)
	}

	_, err = q.PeekByOffset(2)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}
}

func BenchmarkQueueEnqueue(b *testing.B) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		b.Error(err)
	}
	defer q.Drop()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if err := q.Enqueue(NewItemString("value")); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkQueueDequeue(b *testing.B) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		b.Error(err)
	}
	defer q.Drop()

	for n := 0; n < b.N; n++ {
		for i := 1; i <= 10; i++ {
			if err := q.Enqueue(NewItemString("value")); err != nil {
				b.Error(err)
			}
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = q.Dequeue()
	}
}
