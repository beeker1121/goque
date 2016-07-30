package goque_test

import (
	"fmt"

	"github.com/beeker1121/goque"
)

// ExampleQueue demonstrates the implementation of a Goque queue.
func Example_queue() {
	// Open/create a queue.
	q, err := goque.OpenQueue("data_dir")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer q.Close()

	// Enqueue an item.
	item, err := q.Enqueue([]byte("item value"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)         // 1
	fmt.Println(item.Key)        // [0 0 0 0 0 0 0 1]
	fmt.Println(item.Value)      // [105 116 101 109 32 118 97 108 117 101]
	fmt.Println(item.ToString()) // item value

	// Change the item value in the queue.
	item, err = q.Update(item.ID, []byte("new item value"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ToString()) // new item value

	// Dequeue the next item.
	deqItem, err := q.Dequeue()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(deqItem.ToString()) // new item value

	// Delete the queue and its database.
	q.Drop()
}
