package goque

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPrefixQueueClose(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	if _, err = pq.EnqueueString("prefix", "value"); err != nil {
		t.Error(err)
	}

	if pq.Length() != 1 {
		t.Errorf("Expected queue length of 1, got %d", pq.Length())
	}

	pq.Close()

	if _, err = pq.DequeueString("prefix"); err != ErrDBClosed {
		t.Errorf("Expected to get database closed error, got %s", err.Error())
	}

	if pq.Length() != 0 {
		t.Errorf("Expected queue length of 0, got %d", pq.Length())
	}
}

func TestPrefixQueueDrop(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
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

func TestPrefixQueueIncompatibleType(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	prq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer prq.Drop()
	prq.Close()

	if _, err = OpenPrefixQueue(file); err != ErrIncompatibleType {
		t.Error("Expected priority queue to return ErrIncompatibleTypes when opening goquePriorityQueue")
	}
}

func TestPrefixQueueEnqueue(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for i := 1; i <= 10; i++ {
		if _, err = pq.EnqueueString("prefix", fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	if pq.Length() != 10 {
		t.Errorf("Expected queue size of 10, got %d", pq.Length())
	}
}

func TestPrefixQueueDequeue(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for i := 1; i <= 10; i++ {
		if _, err = pq.EnqueueString("prefix", fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	if pq.Length() != 10 {
		t.Errorf("Expected queue length of 10, got %d", pq.Length())
	}

	deqItem, err := pq.DequeueString("prefix")
	if err != nil {
		t.Error(err)
	}

	if pq.Length() != 9 {
		t.Errorf("Expected queue length of 9, got %d", pq.Length())
	}

	compStr := "value for item 1"

	if deqItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, deqItem.ToString())
	}
}

func TestPrefixQueueEncodeDecodePointerJSON(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
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

	if _, err = pq.EnqueueObjectAsJSON([]byte("prefix"), obj); err != nil {
		t.Error(err)
	}

	item, err := pq.Dequeue([]byte("prefix"))
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

func TestPrefixQueuePeek(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	compStr := "value for item"

	if _, err = pq.EnqueueString("prefix", compStr); err != nil {
		t.Error(err)
	}

	peekItem, err := pq.PeekString("prefix")
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	if pq.Length() != 1 {
		t.Errorf("Expected queue length of 1, got %d", pq.Length())
	}
}

func TestPrefixQueuePeekByID(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for i := 1; i <= 10; i++ {
		if _, err = pq.EnqueueString("prefix", fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	compStr := "value for item 3"

	peekItem, err := pq.PeekByIDString("prefix", 3)
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	if pq.Length() != 10 {
		t.Errorf("Expected queue length of 10, got %d", pq.Length())
	}
}

func TestPrefixQueueUpdate(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for i := 1; i <= 10; i++ {
		if _, err = pq.EnqueueString("prefix", fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	item, err := pq.PeekByIDString("prefix", 3)
	if err != nil {
		t.Error(err)
	}

	oldCompStr := "value for item 3"
	newCompStr := "new value for item 3"

	if item.ToString() != oldCompStr {
		t.Errorf("Expected string to be '%s', got '%s'", oldCompStr, item.ToString())
	}

	updatedItem, err := pq.UpdateString("prefix", item.ID, newCompStr)
	if err != nil {
		t.Error(err)
	}

	if updatedItem.ToString() != newCompStr {
		t.Errorf("Expected current item value to be '%s', got '%s'", newCompStr, item.ToString())
	}

	newItem, err := pq.PeekByIDString("prefix", 3)
	if err != nil {
		t.Error(err)
	}

	if newItem.ToString() != newCompStr {
		t.Errorf("Expected new item value to be '%s', got '%s'", newCompStr, item.ToString())
	}
}

func TestPrefixQueueUpdateString(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for i := 1; i <= 10; i++ {
		if _, err = pq.EnqueueString("prefix", fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	item, err := pq.PeekByIDString("prefix", 3)
	if err != nil {
		t.Error(err)
	}

	oldCompStr := "value for item 3"
	newCompStr := "new value for item 3"

	if item.ToString() != oldCompStr {
		t.Errorf("Expected string to be '%s', got '%s'", oldCompStr, item.ToString())
	}

	updatedItem, err := pq.UpdateString("prefix", item.ID, newCompStr)
	if err != nil {
		t.Error(err)
	}

	if updatedItem.ToString() != newCompStr {
		t.Errorf("Expected current item value to be '%s', got '%s'", newCompStr, item.ToString())
	}

	newItem, err := pq.PeekByIDString("prefix", 3)
	if err != nil {
		t.Error(err)
	}

	if newItem.ToString() != newCompStr {
		t.Errorf("Expected new item value to be '%s', got '%s'", newCompStr, item.ToString())
	}
}

func TestPrefixQueueUpdateObject(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	type object struct {
		Value int
	}

	for i := 1; i <= 10; i++ {
		if _, err = pq.EnqueueObject([]byte("prefix"), object{i}); err != nil {
			t.Error(err)
		}
	}

	item, err := pq.PeekByIDString("prefix", 3)
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

	updatedItem, err := pq.UpdateObject([]byte("prefix"), item.ID, newCompObj)
	if err != nil {
		t.Error(err)
	}

	if err := updatedItem.ToObject(&obj); err != nil {
		t.Error(err)
	}

	if obj != newCompObj {
		t.Errorf("Expected current object to be '%+v', got '%+v'", newCompObj, obj)
	}

	newItem, err := pq.PeekByIDString("prefix", 3)
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

func TestPrefixQueueUpdateObjectAsJSON(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
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

	for i := 1; i <= 10; i++ {
		obj := object{
			Value: i,
			SubObject: subObject{
				Value: &i,
			},
		}

		if _, err = pq.EnqueueObjectAsJSON([]byte("prefix"), obj); err != nil {
			t.Error(err)
		}
	}

	item, err := pq.PeekByIDString("prefix", 3)
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

	updatedItem, err := pq.UpdateObjectAsJSON([]byte("prefix"), item.ID, newCompObj)
	if err != nil {
		t.Error(err)
	}

	if err := updatedItem.ToObjectFromJSON(&obj); err != nil {
		t.Error(err)
	}

	if *obj.SubObject.Value != *newCompObj.SubObject.Value {
		t.Errorf("Expected current object subobject value to be '%+v', got '%+v'", *newCompObj.SubObject.Value, *obj.SubObject.Value)
	}

	newItem, err := pq.PeekByIDString("prefix", 3)
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

func TestPrefixQueueUpdateOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	for i := 1; i <= 10; i++ {
		if _, err = pq.EnqueueString("prefix", fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	if pq.Length() != 10 {
		t.Errorf("Expected queue length of 10, got %d", pq.Length())
	}

	deqItem, err := pq.DequeueString("prefix")
	if err != nil {
		t.Error(err)
	}

	if pq.Length() != 9 {
		t.Errorf("Expected queue length of 9, got %d", pq.Length())
	}

	if _, err = pq.Update([]byte("prefix"), deqItem.ID, []byte(`new value`)); err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}

	if _, err = pq.Update([]byte("prefix"), deqItem.ID+1, []byte(`new value`)); err != nil {
		t.Error(err)
	}
}

func TestPrefixQueueEmpty(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.EnqueueString("prefix", "value for item")
	if err != nil {
		t.Error(err)
	}

	_, err = pq.DequeueString("prefix")
	if err != nil {
		t.Error(err)
	}

	_, err = pq.DequeueString("prefix")
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}
}

func TestPrefixQueueOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.EnqueueString("prefix", "value for item")
	if err != nil {
		t.Error(err)
	}

	_, err = pq.PeekByIDString("prefix", 2)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get queue out of bounds error, got %s", err.Error())
	}
}

func TestPrefixQueueRecover(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()

	_, err = pq.EnqueueString("prefix", "value for item")
	if err != nil {
		t.Error(err)
	}

	if err = pq.Close(); err != nil {
		t.Error(err)
	}
	if err = os.Remove(file + "/MANIFEST-000000"); err != nil {
		t.Error(err)
	}

	if pq, err = OpenPrefixQueue(file); !IsCorrupted(err) {
		t.Errorf("Expected corruption error, got %s", err)
	}
	if pq, err = RecoverPrefixQueue(file); err != nil {
		t.Error(err)
	}
}

func BenchmarkPrefixQueueEnqueue(b *testing.B) {
	// Open test database
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		b.Error(err)
	}
	defer pq.Drop()

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, _ = pq.Enqueue([]byte("prefix"), []byte("value"))
	}
}

func BenchmarkPrefixQueueDequeue(b *testing.B) {
	// Open test database
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPrefixQueue(file)
	if err != nil {
		b.Error(err)
	}
	defer pq.Drop()

	// Fill with dummy data
	for n := 0; n < b.N; n++ {
		if _, err = pq.Enqueue([]byte("prefix"), []byte("value")); err != nil {
			b.Error(err)
		}
	}

	// Start benchmark
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, _ = pq.Dequeue([]byte("prefix"))
	}
}
