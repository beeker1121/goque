package goque_test

import (
	"fmt"

	"github.com/beeker1121/goque"
)

// ExampleObject demonstrates enqueuing a struct object.
func Example_object() {
	// Open/create a queue.
	q, err := goque.OpenQueue("data_dir")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer q.Close()

	// Define our struct.
	type object struct {
		X int
		Y int
	}

	// Create a new item with our struct.
	item, err := goque.NewItemObject(object{X: 1, Y: 2})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Enqueue the item.
	err = q.Enqueue(item)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)    // 1
	fmt.Println(item.Key)   // [0 0 0 0 0 0 0 1]
	fmt.Println(item.Value) // [105 116 101 109 32 118 97 108 117 101]

	// Dequeue the item.
	deqItem, err := q.Dequeue()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a variable to hold our decoded struct in.
	var obj object

	// Decode the item into our struct type.
	if err := deqItem.ToObject(&obj); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", obj) // {X:1 Y:2}
}
