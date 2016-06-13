package goque

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPriorityQueueDrop(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}

	if _, err = os.Stat(file); os.IsNotExist(err) {
		t.Error(err)
	}

	pq.Drop()

	if _, err = os.Stat(file); err == nil {
		t.Error("Expected directory for test database to have been deleted")
	}
}

func TestPriorityQueueEnqueue(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			item := NewPriorityItemString(fmt.Sprintf("value for item %d", i), uint8(p))
			if err = pq.Enqueue(item); err != nil {
				t.Error(err)
			}
		}
	}

	if pq.Length() != 50 {
		t.Errorf("Expected queue size of 50, got %d", pq.Length())
	}
}

func TestPriorityQueueDequeueAsc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			item := NewPriorityItemString(fmt.Sprintf("value for item %d", i), uint8(p))
			if err = pq.Enqueue(item); err != nil {
				t.Error(err)
			}
		}
	}

	if pq.Length() != 50 {
		t.Errorf("Expected queue length of 1, got %d", pq.Length())
	}

	deqItem, err := pq.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if pq.Length() != 49 {
		t.Errorf("Expected queue length of 49, got %d", pq.Length())
	}

	compStr := "value for item 1"

	if deqItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", deqItem.Priority)
	}

	if deqItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, deqItem.ToString())
	}
}

func TestPriorityQueueDequeueDesc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, DESC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			item := NewPriorityItemString(fmt.Sprintf("value for item %d", i), uint8(p))
			if err = pq.Enqueue(item); err != nil {
				t.Error(err)
			}
		}
	}

	if pq.Length() != 50 {
		t.Errorf("Expected queue length of 1, got %d", pq.Length())
	}

	deqItem, err := pq.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if pq.Length() != 49 {
		t.Errorf("Expected queue length of 49, got %d", pq.Length())
	}

	compStr := "value for item 1"

	if deqItem.Priority != 4 {
		t.Errorf("Expected priority level to be 4, got %d", deqItem.Priority)
	}

	if deqItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, deqItem.ToString())
	}
}

func TestPriorityQueuePeek(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			item := NewPriorityItemString(fmt.Sprintf("value for item %d", i), uint8(p))
			if err = pq.Enqueue(item); err != nil {
				t.Error(err)
			}
		}
	}

	compStr := "value for item 1"

	peekItem, err := pq.Peek()
	if err != nil {
		t.Error(err)
	}

	if peekItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", peekItem.Priority)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	secondPeekItem, err := pq.Peek()
	if err != nil {
		t.Error(err)
	}

	if secondPeekItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", peekItem.Priority)
	}

	if secondPeekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, secondPeekItem.ToString())
	}
}

func TestPriorityQueuePeekByOffsetAsc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			item := NewPriorityItemString(fmt.Sprintf("value for item %d", i), uint8(p))
			if err = pq.Enqueue(item); err != nil {
				t.Error(err)
			}
		}
	}

	compStrFirst := "value for item 1"
	compStrLast := "value for item 10"
	compStr := "value for item 3"

	peekFirstItem, err := pq.PeekByOffset(0)
	if err != nil {
		t.Error(err)
	}

	if peekFirstItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", peekFirstItem.Priority)
	}

	if peekFirstItem.ToString() != compStrFirst {
		t.Errorf("Expected string to be '%s', got '%s'", compStrFirst, peekFirstItem.ToString())
	}

	peekLastItem, err := pq.PeekByOffset(49)
	if err != nil {
		t.Error(err)
	}

	if peekLastItem.Priority != 4 {
		t.Errorf("Expected priority level to be 4, got %d", peekLastItem.Priority)
	}

	if peekLastItem.ToString() != compStrLast {
		t.Errorf("Expected string to be '%s', got '%s'", compStrLast, peekLastItem.ToString())
	}

	peekItem, err := pq.PeekByOffset(22)
	if err != nil {
		t.Error(err)
	}

	if peekItem.Priority != 2 {
		t.Errorf("Expected priority level to be 2, got %d", peekItem.Priority)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	secondPeekItem, err := pq.PeekByOffset(22)
	if err != nil {
		t.Error(err)
	}

	if secondPeekItem.Priority != 2 {
		t.Errorf("Expected priority level to be 2, got %d", secondPeekItem.Priority)
	}

	if secondPeekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, secondPeekItem.ToString())
	}
}

func TestPriorityQueuePeekByOffsetDesc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, DESC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			item := NewPriorityItemString(fmt.Sprintf("value for item %d", i), uint8(p))
			if err = pq.Enqueue(item); err != nil {
				t.Error(err)
			}
		}
	}

	compStrFirst := "value for item 1"
	compStrLast := "value for item 10"
	compStr := "value for item 3"

	peekFirstItem, err := pq.PeekByOffset(0)
	if err != nil {
		t.Error(err)
	}

	if peekFirstItem.Priority != 4 {
		t.Errorf("Expected priority level to be 4, got %d", peekFirstItem.Priority)
	}

	if peekFirstItem.ToString() != compStrFirst {
		t.Errorf("Expected string to be '%s', got '%s'", compStrFirst, peekFirstItem.ToString())
	}

	peekLastItem, err := pq.PeekByOffset(49)
	if err != nil {
		t.Error(err)
	}

	if peekLastItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", peekLastItem.Priority)
	}

	if peekLastItem.ToString() != compStrLast {
		t.Errorf("Expected string to be '%s', got '%s'", compStrLast, peekLastItem.ToString())
	}

	peekItem, err := pq.PeekByOffset(32)
	if err != nil {
		t.Error(err)
	}

	if peekItem.Priority != 1 {
		t.Errorf("Expected priority level to be 0, got %d", peekItem.Priority)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	secondPeekItem, err := pq.PeekByOffset(32)
	if err != nil {
		t.Error(err)
	}

	if secondPeekItem.Priority != 1 {
		t.Errorf("Expected priority level to be 1, got %d", secondPeekItem.Priority)
	}

	if secondPeekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, secondPeekItem.ToString())
	}
}

func TestPriorityQueuePeekByPriorityID(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			item := NewPriorityItemString(fmt.Sprintf("value for item %d", i), uint8(p))
			if err = pq.Enqueue(item); err != nil {
				t.Error(err)
			}
		}
	}

	compStr := "value for item 3"

	peekItem, err := pq.PeekByPriorityID(1, 3)
	if err != nil {
		t.Error(err)
	}

	if peekItem.Priority != 1 {
		t.Errorf("Expected priority level to be 1, got %d", peekItem.Priority)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	secondPeekItem, err := pq.PeekByPriorityID(1, 3)
	if err != nil {
		t.Error(err)
	}

	if secondPeekItem.Priority != 1 {
		t.Errorf("Expected priority level to be 1, got %d", secondPeekItem.Priority)
	}

	if secondPeekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, secondPeekItem.ToString())
	}
}
