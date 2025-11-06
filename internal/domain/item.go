package domain

import (
	"time"
)

// Item represents a business entity in the domain
type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ItemEvent represents a domain event for Kafka
type ItemEvent struct {
	EventType string    `json:"event_type"`
	ItemID    string    `json:"item_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      Item      `json:"data"`
}

// Validate validates the item
func (i *Item) Validate() error {
	if i.Name == "" {
		return ErrInvalidItem
	}
	if i.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}
