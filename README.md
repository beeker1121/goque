# Goque [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/beeker1121/goque) [![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/beeker1121/goque/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/beeker1121/goque)](https://goreportcard.com/report/github.com/beeker1121/goque) [![Build Status](https://travis-ci.org/beeker1121/goque.svg?branch=master)](https://travis-ci.org/beeker1121/goque)

Goque provides embedded, disk-based implementations of stack, queue, and priority queue data structures.

Motivation for creating this project was the need for a persistent priority queue that remained performant while growing well beyond the available memory of a given machine. While there are many packages for Go offering queues, they all seem to be memory based and/or standalone solutions that are not embeddable within an application.

Instead of using an in-memory heap structure to store data, everything is stored using the [Go port of LevelDB](https://github.com/syndtr/goleveldb). This results in very little memory being used no matter the size of the database, while read and write performance remains near constant.

**IMPORTANT UPDATE: Goque has been updated to v2, which introduces a completely new API for the Enqueue, Push, and Update methods. Please refer the [v2 release](https://github.com/beeker1121/goque/releases/tag/v2.0.0) for more information on exactly what has been changed. The prior API can still be found at [release v1.0.2](https://github.com/beeker1121/goque/releases/tag/v1.0.2).**

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

Push an item:

```go
item, err := s.Push([]byte("item value"))
// or
item, err := s.PushString("item value")
// or
item, err := s.PushObject(Object{X:1})
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
item, err := s.Update(1, []byte("new value"))
// or
item, err := s.UpdateString(1, "new value")
// or
item, err := s.UpdateObject(1, Object{X:2})
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

Enqueue an item:

```go
item, err := q.Enqueue([]byte("item value"))
// or
item, err := q.EnqueueString("item value")
// or
item, err := q.EnqueueObject(Object{X:1})
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
item, err := q.Update(1, []byte("new value"))
// or
item, err := q.UpdateString(1, "new value")
// or
item, err := q.UpdateObject(1, Object{X:2})
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

Enqueue an item:

```go
item, err := pq.Enqueue(0, []byte("item value"))
// or
item, err := pq.EnqueueString(0, "item value")
// or
item, err := pq.EnqueueObject(0, Object{X:1})
```

Dequeue an item:

```go
item, err := pq.Dequeue()
// or
item, err := pq.DequeueByPriority(0)
...
fmt.Println(item.ID)         // 1
fmt.Println(item.Priority)   // 0
fmt.Println(item.Key)        // [0 58 0 0 0 0 0 0 0 1]
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
item, err := pq.Update(0, 1, []byte("new value"))
// or
item, err := pq.UpdateString(0, 1, "new value")
// or
item, err := pq.UpdateObject(0, 1, Object{X:2})
```

Delete the priority queue and underlying database:

```go
pq.Drop()
```

## Benchmarks

Benchmarks were run on a Google Compute Engine n1-standard-1 machine (1 vCPU 3.75 GB of RAM):

```
$ go test -bench=.
PASS
BenchmarkPriorityQueueEnqueue     200000              8104 ns/op             522 B/op          7 allocs/op
BenchmarkPriorityQueueDequeue     200000             18622 ns/op            1166 B/op         17 allocs/op
BenchmarkQueueEnqueue             200000              8049 ns/op             487 B/op          7 allocs/op
BenchmarkQueueDequeue             200000             18970 ns/op            1089 B/op         17 allocs/op
BenchmarkStackPush                200000              8145 ns/op             487 B/op          7 allocs/op
BenchmarkStackPop                 200000             18947 ns/op            1097 B/op         17 allocs/op
ok      github.com/beeker1121/goque     22.549s
```

## Thanks

**syndtr** ([https://github.com/syndtr](https://github.com/syndtr)) - LevelDB port to Go  
**bogdanovich** ([https://github.com/bogdanovich/siberite](https://github.com/bogdanovich/siberite)) - Server based queue for Go using LevelDB  
**connor4312** ([https://github.com/connor4312](https://github.com/connor4312)) - Recommending BoltDB/LevelDB, helping with structure  
**bwmarrin** ([https://github.com/bwmarrin](https://github.com/bwmarrin)) - Recommending BoltDB/LevelDB  
**zeroZshadow** ([https://github.com/zeroZshadow](https://github.com/zeroZshadow)) - Code review and optimization  
**nstafie** ([https://github.com/nstafie](https://github.com/nstafie)) - Help with structure