package goque

import (
	"errors"
)

var (
	// ErrIncompatibleType is returned when trying to open stored
	// Goque data structure type is incompatible with the opener.
	ErrIncompatibleType = errors.New("goque: Opener type is incompatible with stored type")

	// ErrEmpty is returned when the queue is empty.
	ErrEmpty = errors.New("goque: The queue is empty")

	// ErrOutOfBounds is returned when the ID used to lookup an item
	// in the queue is outside the current range of the queue.
	ErrOutOfBounds = errors.New("goque: ID used is out of the range of the queue")

	// ErrDBClosed is returned when the Close function has already
	// been called, causing the stack, queue, or priority queue to
	// close, as well as its underlying database.
	ErrDBClosed = errors.New("goque: The database is closed")
)
