package goque_test

import (
	"fmt"

	"github.com/beeker1121/goque"
)

// ExampleQueue demonstrates the implementation of a Goque queue.
func Example_priorityQueue() {
	// Open/create a priority queue.
	pq, err := goque.OpenPriorityQueue("data_dir", goque.ASC)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pq.Close()

	// Create a new item.
	item := goque.NewPriorityItem([]byte("item value"), 0)

	// Enqueue the item.
	err = pq.Enqueue(item)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)         // 1
	fmt.Println(item.Priority)   // 0
	fmt.Println(item.Key)        // [0 0 0 0 0 0 0 1]
	fmt.Println(item.Value)      // [105 116 101 109 32 118 97 108 117 101]
	fmt.Println(item.ToString()) // item value

	// Change the item value in the queue.
	err = pq.Update(item, []byte("new item value"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ToString()) // new item value

	// Dequeue the next item.
	deqItem, err := pq.Dequeue()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(deqItem.ToString()) // new item value
}
