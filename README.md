# Goque [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/beeker1121/goque) [![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/beeker1121/goque/master/LICENSE) [![GoReport](http://img.shields.io/badge/go_report-A+-brightgreen.svg)](https://goreportcard.com/report/github.com/beeker1121/goque)

Goque provides embedded, disk-based implementations of stack, queue, and priority queue data structures.

Motivation for creating this project was the need for a persistent priority queue that remained performant while growing well beyond the available memory of a given machine. While there are many packages for Go offering queues, they all seem to be memory based and/or standalone solutions that are not embeddable within an application.

The Go implementation of LevelDB is used as the backend for stacks and queues.

## Features

- Provides stack (LIFO), queue (FIFO), and priority queue structures.
- Persistent, disk-based.
- Optimized for fast inserts and reads.
- Designed to work with large datasets outside of RAM/memory.

## Installation

Import to your project:

```go
import "github.com/beeker1121/goque"
```

Fetch the package from GitHub:

```sh
go get github.com/beeker1121/goque
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
```

Push an item:

```go
err := s.Push(item)
...
```

Pop an item:

```go
item, err := s.Pop()
...
fmt.Println(item.ID) // 1
fmt.Println(item.Key) // [0 0 0 0 0 0 0 1]
fmt.Println(item.Value) // [105 116 101 109 32 118 97 108 117 101]
fmt.Println(item.ToString) // item value
```

Peek the next stack item:

```go
item, err := s.Peek()
...
item, err := s.PeekByOffset(1)
...
item, err := s.PeekByID(1)
...
```

Update an item in the stack:

```go
err := s.Update(item, []byte("new value"))
...
// or
err := s.UpdateString(item, "new value")
...
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
```

Enqueue an item:

```go
err := q.Enqueue(item)
...
```

Dequeue an item:

```go
item, err := q.Dequeue()
...
fmt.Println(item.ID) // 1
fmt.Println(item.Key) // [0 0 0 0 0 0 0 1]
fmt.Println(item.Value) // [105 116 101 109 32 118 97 108 117 101]
fmt.Println(item.ToString) // item value
```

Peek the next queue item:

```go
item, err := q.Peek()
...
item, err := q.PeekByOffset(1)
...
item, err := q.PeekByID(1)
...
```

Update an item in the queue:

```go
err := q.Update(item, []byte("new value"))
...
// or
err := q.UpdateString(item, "new value")
...
```

Delete the queue and underlying database:

```go
q.Drop()
```

## Thanks

**syndtr** ([https://github.com/syndtr](https://github.com/syndtr)) - LevelDB port to Go  
**connor4312** ([https://github.com/connor4312](https://github.com/connor4312)) - Recommending BoltDB/LevelDB, helping with structure  
**bwmarrin** ([https://github.com/bwmarrin](https://github.com/bwmarrin)) - Recommending BoltDB/LevelDB  
**zeroZshadow** ([https://github.com/zeroZshadow](https://github.com/zeroZshadow)) - Code review and optimization  