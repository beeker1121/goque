package goque

import (
	"encoding/binary"
)

// Item represents an entry in either a stack or queue.
type Item struct {
	ID    uint64
	Key   []byte
	Value []byte
}

// NewItem creates a new item for use with a stack or queue.
func NewItem(value []byte) *Item {
	return &Item{Value: value}
}

// NewItemString is a helper function for NewItem that accepts a
// value as a string rather than a byte slice.
func NewItemString(value string) *Item {
	return NewItem([]byte(value))
}

// ToString returns the item value as a string.
func (i *Item) ToString() string {
	return string(i.Value)
}

// PriorityItem represents an entry in a priority queue.
type PriorityItem struct {
	ID       uint64
	Priority uint8
	Key      []byte
	Value    []byte
}

// ToString returns the item value as a string.
//func (pi *PriorityItem) ToString() string {}

// idToKey converts and returns the given ID to a key.
func idToKey(id uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, id)
	return key
}

// keyToID converts and returns the given key to an ID.
func keyToID(key []byte) uint64 {
	return binary.BigEndian.Uint64(key)
}
