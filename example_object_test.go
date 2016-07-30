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

	// Enqueue an object.
	item, err := q.EnqueueObject(object{X: 1, Y: 2})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)  // 1
	fmt.Println(item.Key) // [0 0 0 0 0 0 0 1]

	// Dequeue an item.
	deqItem, err := q.Dequeue()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create variable to hold our object in.
	var obj object

	// Decode item into our struct type.
	if err := deqItem.ToObject(&obj); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", obj) // {X:1 Y:2}

	// Delete the queue and its database.
	q.Drop()
}
