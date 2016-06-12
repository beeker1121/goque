package goque_test

import (
	"fmt"

	"goque"
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

	// Create a new item.
	item := goque.NewItem([]byte("item value"))

	// Enqueue the item.
	err = q.Enqueue(item)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)         // 1
	fmt.Println(item.Key)        // [0 0 0 0 0 0 0 1]
	fmt.Println(item.Value)      // [105 116 101 109 32 118 97 108 117 101]
	fmt.Println(item.ToString()) // item value

	// Change the item value in the queue.
	err = q.Update(item, []byte("new item value"))
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
}
