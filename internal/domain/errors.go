package domain

import "errors"

var (
	// ErrItemNotFound is returned when an item is not found
	ErrItemNotFound = errors.New("item not found")
	
	// ErrInvalidItem is returned when an item is invalid
	ErrInvalidItem = errors.New("invalid item: name is required")
	
	// ErrInvalidPrice is returned when price is invalid
	ErrInvalidPrice = errors.New("invalid price: must be non-negative")
)
