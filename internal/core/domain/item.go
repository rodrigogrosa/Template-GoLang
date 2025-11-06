package domain

import (
	"time"
)

// Item represents a domain entity
type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate performs domain validation
func (i *Item) Validate() error {
	if i.Name == "" {
		return ErrInvalidItemName
	}
	return nil
}
