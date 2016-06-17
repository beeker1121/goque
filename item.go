package goque

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
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

// NewItemObject is a helper function for NewItem that accepts any
// value which it'll be marshalled using encoding/gob.
func NewItemObject(value interface{}) (*Item, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(value); err != nil {
		return nil, err
	}
	return NewItem(buffer.Bytes()), nil
}

// ToString returns the item value as a string.
func (i *Item) ToString() string {
	return string(i.Value)
}

// Unmarshall unmarshalls the item value using encoding/gob.
func (i *Item) Unmarshall(value interface{}) error {
	buffer := bytes.NewBuffer(i.Value)
	dec := gob.NewDecoder(buffer)
	return dec.Decode(value)
}

// PriorityItem represents an entry in a priority queue.
type PriorityItem struct {
	ID       uint64
	Priority uint8
	Key      []byte
	Value    []byte
}

// NewPriorityItem creates a new item for use with a priority queue.
func NewPriorityItem(value []byte, priority uint8) *PriorityItem {
	return &PriorityItem{Priority: priority, Value: value}
}

// NewPriorityItemString is a helper function for NewPriorityItem
// that accepts a value as a string rather than a byte slice.
func NewPriorityItemString(value string, priority uint8) *PriorityItem {
	return NewPriorityItem([]byte(value), priority)
}

// NewPriorityItemObject is a helper function for NewPriorityItem
// that accepts any value which it'll be marshalled using encoding/gob.
func NewPriorityItemObject(value interface{}, priority uint8) (*PriorityItem, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(value); err != nil {
		return nil, err
	}
	return NewPriorityItem(buffer.Bytes(), priority), nil
}

// ToString returns the priority item value as a string.
func (pi *PriorityItem) ToString() string {
	return string(pi.Value)
}

// Unmarshall unmarshalls the item value using encoding/gob.
func (pi *PriorityItem) Unmarshall(value interface{}) error {
	buffer := bytes.NewBuffer(pi.Value)
	dec := gob.NewDecoder(buffer)
	return dec.Decode(value)
}

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
