package repository

import (
	"context"
	"sync"

	"github.com/rodrigogrosa/Template-GoLang/internal/core/domain"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/ports"
)

type inMemoryItemRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Item
}

// NewInMemoryItemRepository creates a new in-memory repository
func NewInMemoryItemRepository() ports.ItemRepository {
	return &inMemoryItemRepository{
		items: make(map[string]*domain.Item),
	}
}

func (r *inMemoryItemRepository) Create(ctx context.Context, item *domain.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; exists {
		return domain.ErrItemAlreadyExists
	}

	r.items[item.ID] = item
	return nil
}

func (r *inMemoryItemRepository) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, domain.ErrItemNotFound
	}

	return item, nil
}

func (r *inMemoryItemRepository) List(ctx context.Context, limit, offset int) ([]*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*domain.Item, 0)
	count := 0

	for _, item := range r.items {
		if count < offset {
			count++
			continue
		}
		if len(items) >= limit {
			break
		}
		items = append(items, item)
		count++
	}

	return items, nil
}

func (r *inMemoryItemRepository) Update(ctx context.Context, item *domain.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; !exists {
		return domain.ErrItemNotFound
	}

	r.items[item.ID] = item
	return nil
}

func (r *inMemoryItemRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return domain.ErrItemNotFound
	}

	delete(r.items, id)
	return nil
}
