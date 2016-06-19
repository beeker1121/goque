# Goque [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/beeker1121/goque) [![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/beeker1121/goque/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/beeker1121/goque)](https://goreportcard.com/report/github.com/beeker1121/goque) [![Build Status](https://travis-ci.org/beeker1121/goque.svg?branch=master)](https://travis-ci.org/beeker1121/goque)

Goque provides embedded, disk-based implementations of stack, queue, and priority queue data structures.

Motivation for creating this project was the need for a persistent priority queue that remained performant while growing well beyond the available memory of a given machine. While there are many packages for Go offering queues, they all seem to be memory based and/or standalone solutions that are not embeddable within an application.

Instead of using an in-memory heap structure to store data, everything is stored using the [Go port of LevelDB](https://github.com/syndtr/goleveldb). This results in very little memory being used no matter the size of the database, while read and write performance remains near constant.

## Features

- Provides stack (LIFO), queue (FIFO), and priority queue structures.
- Stacks and queues (but not priority queues) are interchangeable.
- Persistent, disk-based.
- Optimized for fast inserts and reads.
- Goroutine safe.
- Designed to work with large datasets outside of RAM/memory.

## Installation

Fetch the package from GitHub:

```sh
go get github.com/beeker1121/goque
```

Import to your project:

```go
import "github.com/beeker1121/goque"
```

## Usage

### Stack

Stack is a LIFO (last in, first out) data structure.

Create or open a stack:

```go
s, err := goque.OpenStack("data_dir")
...
defer s.Close()
```

Create a new item:

```go
item := goque.NewItem([]byte("item value"))
// or
item := goque.NewItemString("item value")
// or
item, err := goque.NewItemObject(Object{X:1})
```

Push an item:

```go
err := s.Push(item)
```

Pop an item:

```go
item, err := s.Pop()
...
fmt.Println(item.ID)         // 1
fmt.Println(item.Key)        // [0 0 0 0 0 0 0 1]
fmt.Println(item.Value)      // [105 116 101 109 32 118 97 108 117 101]
fmt.Println(item.ToString()) // item value

// Decode to object.
var obj Object
err := item.ToObject(&obj)
...
fmt.Printf("%+v\n", obj) // {X:1}
```

Peek the next stack item:

```go
item, err := s.Peek()
// or
item, err := s.PeekByOffset(1)
// or
item, err := s.PeekByID(1)
```

Update an item in the stack:

```go
err := s.Update(item, []byte("new value"))
// or
err := s.UpdateString(item, "new value")
// or
err := s.UpdateObject(item, Object{X:2})
```

Delete the stack and underlying database:

```go
s.Drop()
```

### Queue

Queue is a FIFO (first in, first out) data structure.

#### Methods

Create or open a queue:

```go
q, err := goque.OpenQueue("data_dir")
...
defer q.Close()
```

Create a new item:

```go
item := goque.NewItem([]byte("item value"))
// or
item := goque.NewItemString("item value")
// or
item, err := goque.NewItemObject(Object{X:1})
```

Enqueue an item:

```go
err := q.Enqueue(item)
```

Dequeue an item:

```go
item, err := q.Dequeue()
...
fmt.Println(item.ID)         // 1
fmt.Println(item.Key)        // [0 0 0 0 0 0 0 1]
fmt.Println(item.Value)      // [105 116 101 109 32 118 97 108 117 101]
fmt.Println(item.ToString()) // item value

// Decode to object.
var obj Object
err := item.ToObject(&obj)
...
fmt.Printf("%+v\n", obj) // {X:1}
```

Peek the next queue item:

```go
item, err := q.Peek()
// or
item, err := q.PeekByOffset(1)
// or
item, err := q.PeekByID(1)
```

Update an item in the queue:

```go
err := q.Update(item, []byte("new value"))
// or
err := q.UpdateString(item, "new value")
// or
err := q.UpdateObject(item, Object{X:2})
```

Delete the queue and underlying database:

```go
q.Drop()
```

### Priority Queue

PriorityQueue is a FIFO (first in, first out) queue with priority levels.

#### Methods

Create or open a priority queue:

```go
pq, err := goque.OpenPriorityQueue("data_dir", goque.ASC)
...
defer pq.Close()
```

Create a new item:

```go
item := goque.NewPriorityItem([]byte("item value"), 0)
// or
item := goque.NewPriorityItemString("item value", 0)
// or
item, err := goque.NewPriorityItemObject(Object{X:1}, 0)
```

Enqueue an item:

```go
err := pq.Enqueue(item)
```

Dequeue an item:

```go
item, err := pq.Dequeue()
// or
item, err := pq.DequeueByPriority(0)
...
fmt.Println(item.ID)         // 1
fmt.Println(item.Priority)   // 0
fmt.Println(item.Key)        // [0 0 0 0 0 0 0 1]
fmt.Println(item.Value)      // [105 116 101 109 32 118 97 108 117 101]
fmt.Println(item.ToString()) // item value

// Decode to object.
var obj Object
err := item.ToObject(&obj)
...
fmt.Printf("%+v\n", obj) // {X:1}
```

Peek the next priority queue item:

```go
item, err := pq.Peek()
// or
item, err := pq.PeekByOffset(1)
// or
item, err := pq.PeekByPriorityID(0, 1)
```

Update an item in the priority queue:

```go
err := pq.Update(item, []byte("new value"))
// or
err := pq.UpdateString(item, "new value")
// or
err := pq.UpdateObject(item, Object{X:2})
```

Delete the priority queue and underlying database:

```go
pq.Drop()
```

## Benchmarks

Benchmarks were run on a Google Compute Engine n1-standard-1 machine (1 vCPU 3.75 GB of RAM):

```
go test -bench=.
PASS
BenchmarkPriorityQueueEnqueue     200000              8102 ns/op             442 B/op          5 allocs/op
BenchmarkPriorityQueueDequeue     200000             18602 ns/op            1161 B/op         17 allocs/op
BenchmarkQueueEnqueue             200000              7582 ns/op             399 B/op          5 allocs/op
BenchmarkQueueDequeue             200000             19317 ns/op            1071 B/op         17 allocs/op
BenchmarkStackPush                200000              7847 ns/op             399 B/op          5 allocs/op
BenchmarkStackPop                 200000             18950 ns/op            1081 B/op         17 allocs/op
```

## Thanks

**syndtr** ([https://github.com/syndtr](https://github.com/syndtr)) - LevelDB port to Go  
**bogdanovich** ([https://github.com/bogdanovich/siberite](https://github.com/bogdanovich/siberite)) - Server based queue for Go using LevelDB  
**connor4312** ([https://github.com/connor4312](https://github.com/connor4312)) - Recommending BoltDB/LevelDB, helping with structure  
**bwmarrin** ([https://github.com/bwmarrin](https://github.com/bwmarrin)) - Recommending BoltDB/LevelDB  
**zeroZshadow** ([https://github.com/zeroZshadow](https://github.com/zeroZshadow)) - Code review and optimization  