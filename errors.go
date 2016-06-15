package goque

import (
	"errors"
)

var (
	// ErrIncompatibleType is returned when the opener type is
	// incompatible with the stored Goque type.
	ErrIncompatibleType = errors.New("goque: Opener type is incompatible with stored Goque type")

	// ErrEmpty is returned when the stack or queue is empty.
	ErrEmpty = errors.New("goque: The stack or queue is empty")

	// ErrOutOfBounds is returned when the ID used to lookup an item
	// is outside the range of the stack or queue.
	ErrOutOfBounds = errors.New("goque: ID used is outside the range of the stack or queue")

	// ErrDBClosed is returned when the Close function has already
	// been called, causing the stack or queue to close, as well as
	// its underlying database.
	ErrDBClosed = errors.New("goque: The database is closed")
)
