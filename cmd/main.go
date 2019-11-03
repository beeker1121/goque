package main

import (
	"fmt"
	"log"

	"github.com/beeker1121/goque"
)

type Obj struct {
	A *int
	B int
}

func (o Obj) String() string {
	if o.A != nil {
		return fmt.Sprintf("%v %v %v", o.A, *o.A, o.B)
	} else {
		return fmt.Sprintf("%v %v", o.A, o.B)
	}
}

func main() {
	queue, err := goque.OpenPriorityQueue("queuetest", goque.ASC)
	if err != nil {
		log.Fatal("OpenPriorityQueue()", err)
	}

	o1 := Obj{new(int), 1}
	*o1.A = 2 // not necessary
	// *o1.A = 10 // this will work
	// o1.A = nil // will work, too

	log.Println("Writing to queue:", o1)

	_, err = queue.EnqueueObject(0, o1)
	if err != nil {
		log.Println("EnqueueObject()", o1, err)
	}

	item, err := queue.Dequeue()
	if err != nil {
		log.Fatal("Dequeue()", err)
	}
	o2 := Obj{}
	err = item.ToObject(&o2)
	if err != nil {
		log.Fatal("ToObject()", err)
	}

	log.Println("Read from queue:", o2)

	queue.Drop()
}
