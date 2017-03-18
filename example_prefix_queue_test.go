package goque_test

import (
	"fmt"

	"github.com/beeker1121/goque"
)

// ExamplePrefixQueue demonstrates the implementation of a Goque queue.
func Example_prefixQueue() {
	// Open/create a prefix queue.
	pq, err := goque.OpenPrefixQueue("data_dir")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pq.Close()

	// Enqueue an item.
	item, err := pq.Enqueue([]byte("prefix"), []byte("item value"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)         // 1
	fmt.Println(item.Key)        // [112 114 101 102 105 120 0 0 0 0 0 0 0 0 1]
	fmt.Println(item.Value)      // [105 116 101 109 32 118 97 108 117 101]
	fmt.Println(item.ToString()) // item value

	// Change the item value in the queue.
	item, err = pq.Update([]byte("prefix"), item.ID, []byte("new item value"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ToString()) // new item value

	// Dequeue the next item.
	deqItem, err := pq.Dequeue([]byte("prefix"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(deqItem.ToString()) // new item value

	// Delete the queue and its database.
	pq.Drop()
}
