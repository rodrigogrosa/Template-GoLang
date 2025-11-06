package domain

import "errors"

var (
	ErrItemNotFound      = errors.New("item not found")
	ErrInvalidItemName   = errors.New("invalid item name: cannot be empty")
	ErrInvalidItemID     = errors.New("invalid item id")
	ErrItemAlreadyExists = errors.New("item already exists")
)
