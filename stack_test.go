package goque

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestStackClose(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	if _, err = s.PushString("value"); err != nil {
		t.Error(err)
	}

	if s.Length() != 1 {
		t.Errorf("Expected stack length of 1, got %d", s.Length())
	}

	s.Close()

	if _, err = s.Pop(); err != ErrDBClosed {
		t.Errorf("Expected to get database closed error, got %s", err.Error())
	}

	if s.Length() != 0 {
		t.Errorf("Expected stack length of 0, got %d", s.Length())
	}
}

func TestStackDrop(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}

	if _, err = os.Stat(file); os.IsNotExist(err) {
		t.Error(err)
	}

	s.Drop()

	if _, err = os.Stat(file); err == nil {
		t.Error("Expected directory for test database to have been deleted")
	}
}

func TestStackIncompatibleType(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	pq, err := OpenPriorityQueue(file, ASC)
	if err != nil {
		t.Error(err)
	}
	defer pq.Drop()
	pq.Close()

	if _, err = OpenStack(file); err != ErrIncompatibleType {
		t.Error("Expected stack to return ErrIncompatibleTypes when opening goquePriorityQueue")
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
		if _, err = s.PushString(fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	if s.Length() != 10 {
		t.Errorf("Expected stack size of 10, got %d", s.Length())
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
		if _, err = s.PushString(fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	if s.Length() != 10 {
		t.Errorf("Expected stack length of 10, got %d", s.Length())
	}

	popItem, err := s.Pop()
	if err != nil {
		t.Error(err)
	}

	if s.Length() != 9 {
		t.Errorf("Expected stack length of 9, got %d", s.Length())
	}

	compStr := "value for item 10"

	if popItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, popItem.ToString())
	}
}

func TestStackPushPopPointerJSON(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	q, err := OpenStack(file)
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

	if _, err = q.PushObjectAsJSON(obj); err != nil {
		t.Error(err)
	}

	item, err := q.Pop()
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

func TestStackPeek(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	compStr := "value for item"

	if _, err = s.PushString(compStr); err != nil {
		t.Error(err)
	}

	peekItem, err := s.Peek()
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	if s.Length() != 1 {
		t.Errorf("Expected stack length of 1, got %d", s.Length())
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
		if _, err = s.PushString(fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	compStrFirst := "value for item 10"
	compStrLast := "value for item 1"
	compStr := "value for item 7"

	peekFirstItem, err := s.PeekByOffset(0)
	if err != nil {
		t.Error(err)
	}

	if peekFirstItem.ToString() != compStrFirst {
		t.Errorf("Expected string to be '%s', got '%s'", compStrFirst, peekFirstItem.ToString())
	}

	peekLastItem, err := s.PeekByOffset(9)
	if err != nil {
		t.Error(err)
	}

	if peekLastItem.ToString() != compStrLast {
		t.Errorf("Expected string to be '%s', got '%s'", compStrLast, peekLastItem.ToString())
	}

	peekItem, err := s.PeekByOffset(3)
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	if s.Length() != 10 {
		t.Errorf("Expected stack length of 10, got %d", s.Length())
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
		if _, err = s.PushString(fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	compStr := "value for item 3"

	peekItem, err := s.PeekByID(3)
	if err != nil {
		t.Error(err)
	}

	if peekItem.ToString() != compStr {
		t.Errorf("Expected string to be '%s', got '%s'", compStr, peekItem.ToString())
	}

	if s.Length() != 10 {
		t.Errorf("Expected stack length of 10, got %d", s.Length())
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
		if _, err = s.PushString(fmt.Sprintf("value for item %d", i)); err != nil {
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

	updatedItem, err := s.Update(item.ID, []byte(newCompStr))
	if err != nil {
		t.Error(err)
	}

	if updatedItem.ToString() != newCompStr {
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
		if _, err = s.PushString(fmt.Sprintf("value for item %d", i)); err != nil {
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

	updatedItem, err := s.UpdateString(item.ID, newCompStr)
	if err != nil {
		t.Error(err)
	}

	if updatedItem.ToString() != newCompStr {
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

func TestStackUpdateObject(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	type object struct {
		Value int
	}

	for i := 1; i <= 10; i++ {
		if _, err = s.PushObject(object{i}); err != nil {
			t.Error(err)
		}
	}

	item, err := s.PeekByID(3)
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

	updatedItem, err := s.UpdateObject(item.ID, newCompObj)
	if err != nil {
		t.Error(err)
	}

	if err := updatedItem.ToObject(&obj); err != nil {
		t.Error(err)
	}

	if obj != newCompObj {
		t.Errorf("Expected current object to be '%+v', got '%+v'", newCompObj, obj)
	}

	newItem, err := s.PeekByID(3)
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

func TestStackUpdateObjectAsJSON(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

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

		if _, err = s.PushObjectAsJSON(obj); err != nil {
			t.Error(err)
		}
	}

	item, err := s.PeekByID(3)
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

	updatedItem, err := s.UpdateObjectAsJSON(item.ID, newCompObj)
	if err != nil {
		t.Error(err)
	}

	if err := updatedItem.ToObjectFromJSON(&obj); err != nil {
		t.Error(err)
	}

	if *obj.SubObject.Value != *newCompObj.SubObject.Value {
		t.Errorf("Expected current object subobject value to be '%+v', got '%+v'", *newCompObj.SubObject.Value, *obj.SubObject.Value)
	}

	newItem, err := s.PeekByID(3)
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

func TestStackUpdateOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	for i := 1; i <= 10; i++ {
		if _, err = s.PushString(fmt.Sprintf("value for item %d", i)); err != nil {
			t.Error(err)
		}
	}

	if s.Length() != 10 {
		t.Errorf("Expected stack length of 10, got %d", s.Length())
	}

	popItem, err := s.Pop()
	if err != nil {
		t.Error(err)
	}

	if s.Length() != 9 {
		t.Errorf("Expected stack length of 9, got %d", s.Length())
	}

	if _, err = s.Update(popItem.ID, []byte(`new value`)); err != ErrOutOfBounds {
		t.Errorf("Expected to get stack out of bounds error, got %s", err.Error())
	}

	if _, err = s.Update(popItem.ID-1, []byte(`new value`)); err != nil {
		t.Error(err)
	}
}

func TestStackEmpty(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	_, err = s.PushString("value for item")
	if err != nil {
		t.Error(err)
	}

	_, err = s.Pop()
	if err != nil {
		t.Error(err)
	}

	_, err = s.Pop()
	if err != ErrEmpty {
		t.Errorf("Expected to get empty error, got %s", err.Error())
	}
}

func TestStackOutOfBounds(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	_, err = s.PushString("value for item")
	if err != nil {
		t.Error(err)
	}

	_, err = s.PeekByOffset(2)
	if err != ErrOutOfBounds {
		t.Errorf("Expected to get stack out of bounds error, got %s", err.Error())
	}
}

func TestStackRecover(t *testing.T) {
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		t.Error(err)
	}
	defer s.Drop()

	_, err = s.PushString("value for item")
	if err != nil {
		t.Error(err)
	}

	if err = s.Close(); err != nil {
		t.Error(err)
	}
	if err = os.Remove(file + "/MANIFEST-000000"); err != nil {
		t.Error(err)
	}

	if s, err = OpenStack(file); !IsCorrupted(err) {
		t.Errorf("Expected corruption error, got %s", err)
	}
	if s, err = RecoverStack(file); err != nil {
		t.Error(err)
	}
}

func BenchmarkStackPush(b *testing.B) {
	// Open test database
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		b.Error(err)
	}
	defer s.Drop()

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, _ = s.PushString("value")
	}
}

func BenchmarkStackPop(b *testing.B) {
	// Open test database
	file := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s, err := OpenStack(file)
	if err != nil {
		b.Error(err)
	}
	defer s.Drop()

	// Fill with dummy data
	for n := 0; n < b.N; n++ {
		if _, err = s.PushString("value"); err != nil {
			b.Error(err)
		}
	}

	// Start benchmark
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, _ = s.Pop()
	}
}
