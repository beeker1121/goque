package goque_test

import (
	"fmt"

	"github.com/beeker1121/goque"
)

// ExampleStack demonstrates the implementation of a Goque stack.
func Example_stack() {
	// Open/create a stack.
	s, err := goque.OpenStack("data_dir")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer s.Close()

	// Create a new item.
	item := goque.NewItem([]byte("item value"))

	// Push it onto the stack.
	err = s.Push(item)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)         // 1
	fmt.Println(item.Key)        // [0 0 0 0 0 0 0 1]
	fmt.Println(item.Value)      // [105 116 101 109 32 118 97 108 117 101]
	fmt.Println(item.ToString()) // item value

	// Change the item value in the stack.
	err = s.Update(item, []byte("new item value"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ToString()) // new item value

	// Pop an item off the stack.
	popItem, err := s.Pop()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(popItem.ToString()) // new item value
}
