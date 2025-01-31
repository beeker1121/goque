package goque

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestQueueClose(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	if _, err = q.EnqueueString("value"); err != nil {
		t.Error(err)
	}

	if q.Length() != 1 {
		t.Errorf("Expected queue length of 1, got %d", q.Length())
	}

	q.Close()

	if _, err = q.Dequeue(); err != ErrDBClosed {
		t.Errorf("Expected to get database closed error, got %s", err.Error())
	}

	if q.Length() != 0 {
		t.Errorf("Expected queue length of 0, got %d", q.Length())
	}
}

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

func TestQueueIncompatibleType(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()
	pq.Close()

	if _, err = OpenQueue(file); err != ErrIncompatibleType {
		t.Error("Expected priority queue to return ErrIncompatibleTypes when opening goquePriorityQueue")
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
		if _, err = q.EnqueueString(fmt.Sprintf("value for item %d", i)); err != nil {
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
		if _, err = q.EnqueueString(fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	if q.Length() != 10 {
		t.Errorf("Expected queue length of 10, got %d", q.Length())
	}

	deqItem, err := q.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if q.Length() != 9 {
		t.Errorf("Expected queue length of 9, got %d", q.Length())
	}

	compStr := "value for item 1"

	if deqItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, deqItem.ToString())
	}
}

func TestQueueEncodeDecodePointerJSON(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

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

	if _, err = q.EnqueueObjectAsJSON(obj); err != nil {
		t.Error(err)
	}

	item, err := q.Dequeue()
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

func TestQueuePeek(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	compStr := "value for item"

	if _, err = q.EnqueueString(compStr); err != nil {
		t.Error(err)
	}

	peekItem, err := q.Peek()
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	if q.Length() != 1 {
		t.Errorf("Expected queue length of 1, got %d", q.Length())
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
		if _, err = q.EnqueueString(fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	compStrFirst := "value for item 1"
	compStrLast := "value for item 10"
	compStr := "value for item 4"

	peekFirstItem, err := q.PeekByOffset(0)
	if err != nil {
		t.Error(err)
	}

	if peekFirstItem.ToString() != compStrFirst {
		t.Errorf("Expected string to be '%s', got '%s'", compStrFirst, peekFirstItem.ToString())
	}

	peekLastItem, err := q.PeekByOffset(9)
	if err != nil {
		t.Error(err)
	}

	if peekLastItem.ToString() != compStrLast {
		t.Errorf("Expected string to be '%s', got '%s'", compStrLast, peekLastItem.ToString())
	}

	peekItem, err := q.PeekByOffset(3)
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	if q.Length() != 10 {
		t.Errorf("Expected queue length of 10, got %d", q.Length())
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
		if _, err = q.EnqueueString(fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	compStr := "value for item 3"

	peekItem, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	if q.Length() != 10 {
		t.Errorf("Expected queue length of 10, got %d", q.Length())
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
		if _, err = q.EnqueueString(fmt.Sprintf("value for item %d", i)); err != nil {
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

	updatedItem, err := q.Update(item.ID, []byte(newCompStr))
	if err != nil {
		t.Error(err)
	}

	if updatedItem.ToString() != newCompStr {
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
		if _, err = q.EnqueueString(fmt.Sprintf("value for item %d", i)); err != nil {
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

	updatedItem, err := q.UpdateString(item.ID, newCompStr)
	if err != nil {
		t.Error(err)
	}

	if updatedItem.ToString() != newCompStr {
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

func TestQueueUpdateObject(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	type object struct {
		Value int
	}

	for i := 1; i <= 10; i++ {
		if _, err = q.EnqueueObject(object{i}); err != nil {
			t.Error(err)
		}
	}

	item, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	oldCompObj := object{3}
	newCompObj := object{33}

	var obj object
	if err := item.ToObject(&obj); err != nil {
		t.Error(err)
	}

	if obj != oldCompObj {
		t.Errorf("Expected object to be '%+v', got '%+v'", oldCompObj, obj)
	}

	updatedItem, err := q.UpdateObject(item.ID, newCompObj)
	if err != nil {
		t.Error(err)
	}

	if err := updatedItem.ToObject(&obj); err != nil {
		t.Error(err)
	}

	if obj != newCompObj {
		t.Errorf("Expected current object to be '%+v', got '%+v'", newCompObj, obj)
	}

	newItem, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	if err := newItem.ToObject(&obj); err != nil {
		t.Error(err)
	}

	if obj != newCompObj {
		t.Errorf("Expected new object to be '%+v', got '%+v'", newCompObj, obj)
	}
}

func TestQueueUpdateObjectAsJSON(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	type subObject struct {
		Value *int
	}

	type object struct {
		Value     int
		SubObject subObject
	}

	for i := 1; i <= 10; i++ {
		obj := object{
			Value: i,
			SubObject: subObject{
				Value: &i,
			},
		}

		if _, err = q.EnqueueObjectAsJSON(obj); err != nil {
			t.Error(err)
		}
	}

	item, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	oldCompObjVal := 3
	oldCompObj := object{
		Value: 3,
		SubObject: subObject{
			Value: &oldCompObjVal,
		},
	}
	newCompObjVal := 33
	newCompObj := object{
		Value: 33,
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

	updatedItem, err := q.UpdateObjectAsJSON(item.ID, newCompObj)
	if err != nil {
		t.Error(err)
	}

	if err := updatedItem.ToObjectFromJSON(&obj); err != nil {
		t.Error(err)
	}

	if *obj.SubObject.Value != *newCompObj.SubObject.Value {
		t.Errorf("Expected current object subobject value to be '%+v', got '%+v'", *newCompObj.SubObject.Value, *obj.SubObject.Value)
	}

	newItem, err := q.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	if err := newItem.ToObjectFromJSON(&obj); err != nil {
		t.Error(err)
	}

	if *obj.SubObject.Value != *newCompObj.SubObject.Value {
		t.Errorf("Expected current object subobject value to be '%+v', got '%+v'", *newCompObj.SubObject.Value, *obj.SubObject.Value)
	}
}

func TestQueueUpdateOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	for i := 1; i <= 10; i++ {
		if _, err = q.EnqueueString(fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	if q.Length() != 10 {
		t.Errorf("Expected queue length of 10, got %d", q.Length())
	}

	deqItem, err := q.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if q.Length() != 9 {
		t.Errorf("Expected queue length of 9, got %d", q.Length())
	}

	if _, err = q.Update(deqItem.ID, []byte(`new value`)); err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}

	if _, err = q.Update(deqItem.ID+1, []byte(`new value`)); err != nil {
		t.Error(err)
	}
}

func TestQueueEmpty(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	_, err = q.EnqueueString("value for item")
	if err != nil {
		t.Error(err)
	}

	_, err = q.Dequeue()
	if err != nil {
		t.Error(err)
	}

	_, err = q.Dequeue()
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}
}

func TestQueueOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	_, err = q.EnqueueString("value for item")
	if err != nil {
		t.Error(err)
	}

	_, err = q.PeekByOffset(2)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}
}

func TestQueueRecover(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer q.Drop()

	_, err = q.EnqueueString("value for item")
	if err != nil {
		t.Error(err)
	}

	if err = q.Close(); err != nil {
		t.Error(err)
	}
	if err = os.Remove(file + "/MANIFEST-000000"); err != nil {
		t.Error(err)
	}

	if q, err = OpenQueue(file); !IsCorrupted(err) {
		t.Errorf("Expected corruption error, got %s", err)
	}
	if q, err = RecoverQueue(file); err != nil {
		t.Error(err)
	}
}

func BenchmarkQueueEnqueue(b *testing.B) {
	// Open test database
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		b.Error(err)
	}
	defer q.Drop()

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, _ = q.Enqueue([]byte("value"))
	}
}

func BenchmarkQueueDequeue(b *testing.B) {
	// Open test database
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenQueue(file)
	if err != nil {
		b.Error(err)
	}
	defer q.Drop()

	// Fill with dummy data
	for n := 0; n < b.N; n++ {
		if _, err = q.Enqueue([]byte("value")); err != nil {
			b.Error(err)
		}
	}

	// Start benchmark
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, _ = q.Dequeue()
	}
}
