package goque_test

import (
	"fmt"

	"github.com/beeker1121/goque"
)

// ExamplePriorityQueue demonstrates the implementation of a Goque queue.
func Example_priorityQueue() {
	// Open/create a priority queue.
	pq, err := goque.OpenPriorityQueue("data_dir", goque.ASC)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pq.Close()

	// Enqueue the item.
	item, err := pq.Enqueue(0, []byte("item value"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)         // 1
	fmt.Println(item.Priority)   // 0
	fmt.Println(item.Key)        // [0 58 0 0 0 0 0 0 0 1]
	fmt.Println(item.Value)      // [105 116 101 109 32 118 97 108 117 101]
	fmt.Println(item.ToString()) // item value

	// Change the item value in the queue.
	item, err = pq.Update(item.Priority, item.ID, []byte("new item value"))
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

	// Delete the queue and its database.
	pq.Drop()
}
