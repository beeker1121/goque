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
	DataDir  string
	db       *leveldb.DB
	order    order
	levels   [256]*priorityLevel
	curLevel uint8
	isOpen   bool
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

	// Set item ID and key.
	item.ID = level.tail + 1
	item.Key = generateKey(item.Priority, item.ID)

	// Add it to the priority queue.
	err := pq.db.Put(item.Key, item.Value, nil)
	if err == nil {
		level.tail++

		// If this priority level is more important than the curLevel.
		if pq.cmpAsc(item.Priority) || pq.cmpDesc(item.Priority) {
			pq.curLevel = item.Priority
		}
	}

	return err
}

// Dequeue removes the next item in the priority queue and returns it.
func (pq *PriorityQueue) Dequeue() (*PriorityItem, error) {
	pq.Lock()
	defer pq.Unlock()

	// If the current priority level is empty.
	if pq.levels[pq.curLevel].Length() == 0 {
		// Set starting value for curLevel.
		pq.resetCurrentLevel()

		// Try to get the next priority level.
		for i := 0; i < 255; i++ {
			if (pq.cmpAsc(uint8(i)) || pq.cmpDesc(uint8(i))) && pq.levels[uint8(i)].Length() > 0 {
				pq.curLevel = uint8(i)
			}
		}

		// If still empty, return queue empty error.
		if pq.levels[pq.curLevel].Length() == 0 {
			return nil, ErrEmpty
		}
	}

	// Try to get the next item in the current priority level.
	item, err := pq.getItemByPriorityID(pq.curLevel, pq.levels[pq.curLevel].head+1)
	if err != nil {
		return item, err
	}

	// Remove this item from the priority queue.
	if err = pq.db.Delete(item.Key, nil); err != nil {
		return item, err
	}

	// Increment position.
	pq.levels[pq.curLevel].head++

	return item, nil
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

// cmpAsc returns wehther the given priority level is higher than the
// current priority level based on ascending order.
func (pq *PriorityQueue) cmpAsc(level uint8) bool {
	return pq.order == ASC && level < pq.curLevel
}

// cmpAsc returns wehther the given priority level is higher than the
// current priority level based on descending order.
func (pq *PriorityQueue) cmpDesc(level uint8) bool {
	return pq.order == DESC && level > pq.curLevel
}

// resetCurrentLevel resets the current priority level of the queue
// so the highest level can be found.
func (pq *PriorityQueue) resetCurrentLevel() {
	if pq.order == ASC {
		pq.curLevel = 255
	} else if pq.order == DESC {
		pq.curLevel = 0
	}
}

// getItemByID returns an item, if found, for the given ID.
func (pq *PriorityQueue) getItemByPriorityID(priority uint8, id uint64) (*PriorityItem, error) {
	// Check if empty or out of bounds.
	if pq.levels[priority].Length() < 1 {
		return nil, ErrEmpty
	} else if id <= pq.levels[priority].head || id > pq.levels[priority].tail {
		return nil, ErrOutOfBounds
	}

	var err error

	// Create a new PriorityItem.
	item := &PriorityItem{ID: id, Priority: priority, Key: generateKey(priority, id)}
	item.Value, err = pq.db.Get(item.Key, nil)

	return item, err
}

// init initializes the priority queue data.
func (pq *PriorityQueue) init() error {
	// Set starting value for curLevel.
	pq.resetCurrentLevel()

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

			// Since this priority level has item(s), handle updating curLevel.
			if pq.cmpAsc(uint8(i)) || pq.cmpDesc(uint8(i)) {
				pq.curLevel = uint8(i)
			}
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
	// priority + prefixSep = 1 + 1 = 2
	prefix := make([]byte, 2)
	prefix[0] = byte(level)
	prefix[1] = prefixSep[0]
	return prefix
}

func generateKey(priority uint8, id uint64) []byte {
	// prefix + key = 2 + 8 = 10
	key := make([]byte, 10)
	copy(key[0:2], generatePrefix(priority))
	copy(key[2:], idToKey(id))
	return key
}
