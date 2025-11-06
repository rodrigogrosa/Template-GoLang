package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/domain"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/ports"
)

type itemService struct {
	repo     ports.ItemRepository
	producer ports.MessageProducer
}

// NewItemService creates a new item service
func NewItemService(repo ports.ItemRepository, producer ports.MessageProducer) ports.ItemService {
	return &itemService{
		repo:     repo,
		producer: producer,
	}
}

func (s *itemService) CreateItem(ctx context.Context, item *domain.Item) error {
	// Validate domain rules
	if err := item.Validate(); err != nil {
		return err
	}

	// Set metadata
	if item.ID == "" {
		item.ID = uuid.New().String()
	}
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	// Persist
	if err := s.repo.Create(ctx, item); err != nil {
		return err
	}

	// Publish event (non-blocking, best effort)
	if s.producer != nil {
		// Create a context with timeout for the async operation
		go func() {
			publishCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.producer.PublishItemCreated(publishCtx, item)
		}()
	}

	return nil
}

func (s *itemService) GetItem(ctx context.Context, id string) (*domain.Item, error) {
	if id == "" {
		return nil, domain.ErrInvalidItemID
	}
	return s.repo.GetByID(ctx, id)
}

func (s *itemService) ListItems(ctx context.Context, limit, offset int) ([]*domain.Item, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *itemService) UpdateItem(ctx context.Context, item *domain.Item) error {
	if err := item.Validate(); err != nil {
		return err
	}

	// Check if exists
	existing, err := s.repo.GetByID(ctx, item.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrItemNotFound
	}

	item.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, item); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		go func() {
			publishCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.producer.PublishItemUpdated(publishCtx, item)
		}()
	}

	return nil
}

func (s *itemService) DeleteItem(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidItemID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		go func() {
			publishCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.producer.PublishItemDeleted(publishCtx, id)
		}()
	}

	return nil
}
