package repository

import (
	"context"
	"sync"

	"github.com/rodrigogrosa/Template-GoLang/internal/domain"
)

// InMemoryRepository is an in-memory implementation of ItemRepository
type InMemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Item
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: make(map[string]*domain.Item),
	}
}

// Create creates a new item
func (r *InMemoryRepository) Create(ctx context.Context, item *domain.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[item.ID] = item
	return nil
}

// GetByID retrieves an item by ID
func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, domain.ErrItemNotFound
	}

	// Return a copy to prevent external modifications
	itemCopy := *item
	return &itemCopy, nil
}

// List retrieves a list of items
func (r *InMemoryRepository) List(ctx context.Context, limit, offset int) ([]*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*domain.Item, 0, len(r.items))
	for _, item := range r.items {
		itemCopy := *item
		items = append(items, &itemCopy)
	}

	// Apply offset and limit
	start := offset
	if start > len(items) {
		start = len(items)
	}

	end := start + limit
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil
}

// Update updates an existing item
func (r *InMemoryRepository) Update(ctx context.Context, item *domain.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; !exists {
		return domain.ErrItemNotFound
	}

	r.items[item.ID] = item
	return nil
}

// Delete deletes an item
func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return domain.ErrItemNotFound
	}

	delete(r.items, id)
	return nil
}
