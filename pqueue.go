package goque

import (
	"os"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// prefixSep is the prefix separator for each item key.
var prefixSep []byte = []byte(":")

// pqorder defines the priority ordering of the queue.
type order int

// Defines which priority order to dequeue in.
//
// ASC will use priority level 0 as the most important.
// DESC will use priority level 255 as the most important.
const (
	ASC order = iota
	DESC
)

// priorityLevel holds the head and tail position of a priority
// level within the queue.
type priorityLevel struct {
	head uint64
	tail uint64
}

// Length returns the total number of items in this priority level.
func (pl *priorityLevel) Length() uint64 {
	return pl.tail - pl.head
}

// PriorityQueue is a standard FIFO (first in, first out) queue with
// priority levels.
type PriorityQueue struct {
	sync.RWMutex
	DataDir string
	db      *leveldb.DB
	order   order
	levels  [256]*priorityLevel
	isOpen  bool
}

// OpenPriorityQueue opens a priority queue if one exists at the given
// directory. If one does not already exist, a new priority queue is
// created.
func OpenPriorityQueue(dataDir string, order order) (*PriorityQueue, error) {
	var err error

	// Create a new PriorityQueue.
	pq := &PriorityQueue{
		DataDir: dataDir,
		db:      &leveldb.DB{},
		order:   order,
		isOpen:  false,
	}

	// Open database for the priority queue.
	pq.db, err = leveldb.OpenFile(dataDir, nil)
	if err != nil {
		return pq, err
	}

	// Set isOpen and return.
	pq.isOpen = true
	return pq, pq.init()
}

// Enqueue adds an item to the priority queue.
func (pq *PriorityQueue) Enqueue(item *PriorityItem) error {
	pq.Lock()
	defer pq.Unlock()

	// Get the priorityLevel.
	level := pq.levels[item.Priority]

	// Set item ID.
	item.ID = level.tail + 1

	// Set the item key.
	// priority + prefix + key = 1 + 1 + 8 = 10
	key := make([]byte, 10)
	copy(key[0:2], generatePrefix(item.Priority))
	copy(key[2:], idToKey(item.ID))

	// Add it to the priority queue.
	err := pq.db.Put(item.Key, item.Value, nil)
	if err == nil {
		level.tail++
	}

	return err
}

// Length returns the total number of items in the priority queue.
func (pq *PriorityQueue) Length() uint64 {
	var length uint64 = 0
	for _, v := range pq.levels {
		length += v.Length()
	}

	return length
}

// Close closes the LevelDB database of the priority queue.
func (pq *PriorityQueue) Close() {
	// If queue is already closed.
	if !pq.isOpen {
		return
	}

	pq.db.Close()
	pq.isOpen = false
}

// Drop closes and deletes the LevelDB database of the priority queue.
func (pq *PriorityQueue) Drop() {
	pq.Close()
	os.RemoveAll(pq.DataDir)
}

// init initializes the priority queue data.
func (pq *PriorityQueue) init() error {
	// Loop through each priority level.
	for i := 0; i <= 255; i++ {
		// Create a new LevelDB Iterator for this priority level.
		prefix := generatePrefix(uint8(i))
		iter := pq.db.NewIterator(util.BytesPrefix(prefix), nil)

		// Create a new priorityLevel.
		pl := &priorityLevel{
			head: 0,
			tail: 0,
		}

		// Set priority level head to the first item.
		if iter.First() {
			pl.head = keyToID(iter.Key()[2:]) - 1
		}

		// Set priority level tail to the last item.
		if iter.Last() {
			pl.tail = keyToID(iter.Key()[2:])
		}

		if iter.Error() != nil {
			return iter.Error()
		}

		pq.levels[i] = pl
		iter.Release()
	}

	return nil
}

// generatePrefix creates the key prefix for the given priority level.
func generatePrefix(level uint8) []byte {
	prefix := make([]byte, 2)
	prefix[0] = byte(level)
	prefix[1] = prefixSep[0]
	return prefix
}
