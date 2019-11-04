package goque

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"
)

func TestPriorityQueueClose(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
				t.Error(err)
			}
		}
	}

	if pq.Length() != 50 {
		t.Errorf("Expected queue length of 1, got %d", pq.Length())
	}

	pq.Close()

	if _, err = pq.Dequeue(); err != ErrDBClosed {
		t.Errorf("Expected to get database closed error, got %s", err.Error())
	}

	if pq.Length() != 0 {
		t.Errorf("Expected queue length of 0, got %d", pq.Length())
	}
}

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

func TestPriorityQueueIncompatibleType(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()
	q.Close()

	if _, err = OpenPriorityQueue(file, ASC); err != ErrIncompatibleType {
		t.Error("Expected priority queue to return ErrIncompatibleTypes when opening Queue")
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
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
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
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
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
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
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

func TestPriorityQueueDequeueByPriority(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
				t.Error(err)
			}
		}
	}

	if pq.Length() != 50 {
		t.Errorf("Expected queue length of 50, got %d", pq.Length())
	}

	deqItem, err := pq.DequeueByPriority(3)
	if err != nil {
		t.Error(err)
	}

	if pq.Length() != 49 {
		t.Errorf("Expected queue length of 49, got %d", pq.Length())
	}

	compStr := "value for item 1"

	if deqItem.Priority != 3 {
		t.Errorf("Expected priority level to be 1, got %d", deqItem.Priority)
	}

	if deqItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, deqItem.ToString())
	}
}

func TestPriorityQueueEncodeDecodePointerJSON(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, DESC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	type subObject struct {
		Value *int
	}

	type object struct {
		Value     int
		SubObject subObject
	}

	val := 0
	obj := object{
		Value: 0,
		SubObject: subObject{
			Value: &val,
		},
	}

	if _, err = pq.EnqueueObjectAsJSON(0, obj); err != nil {
		t.Error(err)
	}

	item, err := pq.Dequeue()
	if err != nil {
		t.Error(err)
	}

	var itemObj object
	if err := item.ToObjectFromJSON(&itemObj); err != nil {
		t.Error(err)
	}

	if *itemObj.SubObject.Value != 0 {
		t.Errorf("Expected object subobject value to be '0', got '%v'", *itemObj.SubObject.Value)
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
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
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

	if pq.Length() != 50 {
		t.Errorf("Expected queue length of 50, got %d", pq.Length())
	}
}

func TestPriorityQueuePeekByOffsetEmptyAsc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.PeekByOffset(0)
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}

	if _, err = pq.EnqueueString(0, "value"); err != nil {
		t.Error(err)
	}

	_, err = pq.PeekByOffset(0)
	if err != nil {
		t.Errorf("Expected to get nil error, got %s", err.Error())
	}

	if _, err = pq.Dequeue(); err != nil {
		t.Error(err)
	}

	_, err = pq.PeekByOffset(0)
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}
}

func TestPriorityQueuePeekByOffsetEmptyDesc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, DESC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.PeekByOffset(0)
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}

	if _, err = pq.EnqueueString(0, "value"); err != nil {
		t.Error(err)
	}

	_, err = pq.PeekByOffset(0)
	if err != nil {
		t.Errorf("Expected to get nil error, got %s", err.Error())
	}

	if _, err = pq.Dequeue(); err != nil {
		t.Error(err)
	}

	_, err = pq.PeekByOffset(0)
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}
}

func TestPriorityQueuePeekByOffsetBoundsAsc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.PeekByOffset(0)
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}

	if _, err = pq.EnqueueString(0, "value"); err != nil {
		t.Error(err)
	}

	_, err = pq.PeekByOffset(0)
	if err != nil {
		t.Errorf("Expected to get nil error, got %s", err.Error())
	}

	_, err = pq.PeekByOffset(1)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
				t.Error(err)
			}
		}
	}

	_, err = pq.PeekByOffset(50)
	if err != nil {
		t.Errorf("Expected to get nil error, got %s", err.Error())
	}

	_, err = pq.PeekByOffset(51)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}
}

func TestPriorityQueuePeekByOffsetBoundsDesc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, DESC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.PeekByOffset(0)
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}

	if _, err = pq.EnqueueString(0, "value"); err != nil {
		t.Error(err)
	}

	_, err = pq.PeekByOffset(0)
	if err != nil {
		t.Errorf("Expected to get nil error, got %s", err.Error())
	}

	_, err = pq.PeekByOffset(1)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
				t.Error(err)
			}
		}
	}

	_, err = pq.PeekByOffset(50)
	if err != nil {
		t.Errorf("Expected to get nil error, got %s", err.Error())
	}

	_, err = pq.PeekByOffset(51)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
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
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
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

	if pq.Length() != 50 {
		t.Errorf("Expected queue length of 50, got %d", pq.Length())
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
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
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

	if pq.Length() != 50 {
		t.Errorf("Expected queue length of 50, got %d", pq.Length())
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
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
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

	if pq.Length() != 50 {
		t.Errorf("Expected queue length of 50, got %d", pq.Length())
	}
}

func TestPriorityQueueUpdate(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
				t.Error(err)
			}
		}
	}

	item, err := pq.PeekByPriorityID(0, 3)
	if err != nil {
		t.Error(err)
	}

	oldCompStr := "value for item 3"
	newCompStr := "new value for item 3"

	if item.ToString() != oldCompStr {
		t.Errorf("Expected string to be '%s', got '%s'", oldCompStr, item.ToString())
	}

	updatedItem, err := pq.Update(item.Priority, item.ID, []byte(newCompStr))
	if err != nil {
		t.Error(err)
	}

	if updatedItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", item.Priority)
	}

	if updatedItem.ToString() != newCompStr {
		t.Errorf("Expected current item value to be '%s', got '%s'", newCompStr, item.ToString())
	}

	newItem, err := pq.PeekByPriorityID(0, 3)
	if err != nil {
		t.Error(err)
	}

	if newItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", newItem.Priority)
	}

	if newItem.ToString() != newCompStr {
		t.Errorf("Expected new item value to be '%s', got '%s'", newCompStr, item.ToString())
	}
}

func TestPriorityQueueUpdateString(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
				t.Error(err)
			}
		}
	}

	item, err := pq.PeekByPriorityID(0, 3)
	if err != nil {
		t.Error(err)
	}

	oldCompStr := "value for item 3"
	newCompStr := "new value for item 3"

	if item.ToString() != oldCompStr {
		t.Errorf("Expected string to be '%s', got '%s'", oldCompStr, item.ToString())
	}

	updatedItem, err := pq.UpdateString(item.Priority, item.ID, newCompStr)
	if err != nil {
		t.Error(err)
	}

	if updatedItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", item.Priority)
	}

	if updatedItem.ToString() != newCompStr {
		t.Errorf("Expected current item value to be '%s', got '%s'", newCompStr, item.ToString())
	}

	newItem, err := pq.PeekByPriorityID(0, 3)
	if err != nil {
		t.Error(err)
	}

	if newItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", newItem.Priority)
	}

	if newItem.ToString() != newCompStr {
		t.Errorf("Expected new item value to be '%s', got '%s'", newCompStr, item.ToString())
	}
}

func TestPriorityQueueUpdateObject(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	type object struct {
		Priority uint8
		Value    int
	}

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueObject(uint8(p), object{uint8(p), i}); err != nil {
				t.Error(err)
			}
		}
	}

	item, err := pq.PeekByPriorityID(0, 3)
	if err != nil {
		t.Error(err)
	}

	oldCompObj := object{0, 3}
	newCompObj := object{0, 33}

	var obj object
	if err := item.ToObject(&obj); err != nil {
		t.Error(err)
	}

	if obj != oldCompObj {
		t.Errorf("Expected object to be '%+v', got '%+v'", oldCompObj, obj)
	}

	updatedItem, err := pq.UpdateObject(item.Priority, item.ID, newCompObj)
	if err != nil {
		t.Error(err)
	}

	if updatedItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", item.Priority)
	}

	if err := updatedItem.ToObject(&obj); err != nil {
		t.Error(err)
	}

	if obj != newCompObj {
		t.Errorf("Expected current object to be '%+v', got '%+v'", newCompObj, obj)
	}

	newItem, err := pq.PeekByPriorityID(0, 3)
	if err != nil {
		t.Error(err)
	}

	if newItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", newItem.Priority)
	}

	if err := newItem.ToObject(&obj); err != nil {
		t.Error(err)
	}

	if obj != newCompObj {
		t.Errorf("Expected new object to be '%+v', got '%+v'", newCompObj, obj)
	}
}

func TestPriorityQueueUpdateObjectAsJSON(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	type subObject struct {
		Value *int
	}

	type object struct {
		Priority  uint8
		Value     int
		SubObject subObject
	}

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			obj := object{
				Priority: uint8(p),
				Value:    i,
				SubObject: subObject{
					Value: &i,
				},
			}

			if _, err = pq.EnqueueObjectAsJSON(uint8(p), obj); err != nil {
				t.Error(err)
			}
		}
	}

	item, err := pq.PeekByPriorityID(0, 3)
	if err != nil {
		t.Error(err)
	}

	oldCompObjVal := 3
	oldCompObj := object{
		Priority: 0,
		Value:    3,
		SubObject: subObject{
			Value: &oldCompObjVal,
		},
	}
	newCompObjVal := 33
	newCompObj := object{
		Priority: 0,
		Value:    33,
		SubObject: subObject{
			Value: &newCompObjVal,
		},
	}

	var obj object
	if err := item.ToObjectFromJSON(&obj); err != nil {
		t.Error(err)
	}

	if *obj.SubObject.Value != *oldCompObj.SubObject.Value {
		t.Errorf("Expected object subobject value to be '%+v', got '%+v'", *oldCompObj.SubObject.Value, *obj.SubObject.Value)
	}

	updatedItem, err := pq.UpdateObjectAsJSON(item.Priority, item.ID, newCompObj)
	if err != nil {
		t.Error(err)
	}

	if updatedItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", item.Priority)
	}

	if err := updatedItem.ToObjectFromJSON(&obj); err != nil {
		t.Error(err)
	}

	if *obj.SubObject.Value != *newCompObj.SubObject.Value {
		t.Errorf("Expected current object subobject value to be '%+v', got '%+v'", *newCompObj.SubObject.Value, *obj.SubObject.Value)
	}

	newItem, err := pq.PeekByPriorityID(0, 3)
	if err != nil {
		t.Error(err)
	}

	if newItem.Priority != 0 {
		t.Errorf("Expected priority level to be 0, got %d", newItem.Priority)
	}

	if err := newItem.ToObjectFromJSON(&obj); err != nil {
		t.Error(err)
	}

	if *obj.SubObject.Value != *newCompObj.SubObject.Value {
		t.Errorf("Expected current object subobject value to be '%+v', got '%+v'", *newCompObj.SubObject.Value, *obj.SubObject.Value)
	}
}

func TestPriorityQueueUpdateOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 0; p <= 4; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
				t.Error(err)
			}
		}
	}

	if pq.Length() != 50 {
		t.Errorf("Expected queue length of 50, got %d", pq.Length())
	}

	deqItem, err := pq.DequeueByPriority(3)
	if err != nil {
		t.Error(err)
	}

	if pq.Length() != 49 {
		t.Errorf("Expected queue length of 49, got %d", pq.Length())
	}

	if _, err = pq.Update(deqItem.Priority, deqItem.ID, []byte(`new value`)); err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}

	if _, err = pq.Update(deqItem.Priority, deqItem.ID+1, []byte(`new value`)); err != nil {
		t.Error(err)
	}
}

func TestPriorityQueueHigherPriorityAsc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 5; p <= 9; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
				t.Error(err)
			}
		}
	}

	item, err := pq.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if item.Priority != 5 {
		t.Errorf("Expected priority level to be 5, got %d", item.Priority)
	}

	_, err = pq.EnqueueString(2, "value")
	if err != nil {
		t.Error(err)
	}

	higherItem, err := pq.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if higherItem.Priority != 2 {
		t.Errorf("Expected priority level to be 2, got %d", higherItem.Priority)
	}
}

func TestPriorityQueueHigherPriorityDesc(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, DESC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for p := 5; p <= 9; p++ {
		for i := 1; i <= 10; i++ {
			if _, err = pq.EnqueueString(uint8(p), fmt.Sprintf("value for item %d", i)); err != nil {
				t.Error(err)
			}
		}
	}

	item, err := pq.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if item.Priority != 9 {
		t.Errorf("Expected priority level to be 9, got %d", item.Priority)
	}

	_, err = pq.EnqueueString(12, "value")
	if err != nil {
		t.Error(err)
	}

	higherItem, err := pq.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if higherItem.Priority != 12 {
		t.Errorf("Expected priority level to be 12, got %d", higherItem.Priority)
	}
}

func TestPriorityQueueEmpty(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.EnqueueString(0, "value for item")
	if err != nil {
		t.Error(err)
	}

	_, err = pq.Dequeue()
	if err != nil {
		t.Error(err)
	}

	_, err = pq.Dequeue()
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}
}

func TestPriorityQueueOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.EnqueueString(0, "value for item")
	if err != nil {
		t.Error(err)
	}

	_, err = pq.PeekByOffset(2)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}
}

func TestPriorityQueueRecover(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.EnqueueString(0, "value for item")
	if err != nil {
		t.Error(err)
	}

	if err = pq.Close(); err != nil {
		t.Error(err)
	}
	if err = os.Remove(file + "/MANIFEST-000000"); err != nil {
		t.Error(err)
	}

	if pq, err = OpenPriorityQueue(file, ASC); !IsCorrupted(err) {
		t.Errorf("Expected corruption error, got %s", err)
	}
	if pq, err = RecoverPriorityQueue(file, ASC); err != nil {
		t.Error(err)
	}
}

func BenchmarkPriorityQueueEnqueue(b *testing.B) {
	// Open test database
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		b.Error(err)
	}
	defer pq.Drop()

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, _ = pq.EnqueueString(0, "value")
	}
}

func BenchmarkPriorityQueueDequeue(b *testing.B) {
	// Open test database
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		b.Error(err)
	}
	defer pq.Drop()

	// Fill with dummy data
	for n := 0; n < b.N; n++ {
		if _, err = pq.EnqueueString(uint8(math.Mod(float64(n), 255)), "value"); err != nil {
			b.Error(err)
		}
	}

	// Start benchmark
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, _ = pq.Dequeue()
	}
}
