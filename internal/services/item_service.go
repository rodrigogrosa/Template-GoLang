package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rodrigogrosa/Template-GoLang/internal/domain"
	"github.com/rodrigogrosa/Template-GoLang/internal/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("item-service")

// ItemServiceImpl implements the ItemService interface
type ItemServiceImpl struct {
	repo      ports.ItemRepository
	publisher ports.EventPublisher
}

// NewItemService creates a new item service
func NewItemService(repo ports.ItemRepository, publisher ports.EventPublisher) *ItemServiceImpl {
	return &ItemServiceImpl{
		repo:      repo,
		publisher: publisher,
	}
}

// CreateItem creates a new item
func (s *ItemServiceImpl) CreateItem(ctx context.Context, item *domain.Item) error {
	ctx, span := tracer.Start(ctx, "CreateItem", trace.WithAttributes(
		attribute.String("item.name", item.Name),
	))
	defer span.End()

	// Set timestamps and ID
	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	// Validate
	if err := item.Validate(); err != nil {
		span.RecordError(err)
		return err
	}

	// Persist
	if err := s.repo.Create(ctx, item); err != nil {
		span.RecordError(err)
		return err
	}

	// Publish event
	event := &domain.ItemEvent{
		EventType: "item.created",
		ItemID:    item.ID,
		Timestamp: time.Now(),
		Data:      *item,
	}
	if err := s.publisher.PublishItemEvent(ctx, event); err != nil {
		span.RecordError(err)
		// Log error but don't fail the operation
	}

	return nil
}

// GetItem retrieves an item by ID
func (s *ItemServiceImpl) GetItem(ctx context.Context, id string) (*domain.Item, error) {
	ctx, span := tracer.Start(ctx, "GetItem", trace.WithAttributes(
		attribute.String("item.id", id),
	))
	defer span.End()

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return item, nil
}

// ListItems retrieves a list of items
func (s *ItemServiceImpl) ListItems(ctx context.Context, limit, offset int) ([]*domain.Item, error) {
	ctx, span := tracer.Start(ctx, "ListItems", trace.WithAttributes(
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	))
	defer span.End()

	items, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return items, nil
}

// UpdateItem updates an existing item
func (s *ItemServiceImpl) UpdateItem(ctx context.Context, item *domain.Item) error {
	ctx, span := tracer.Start(ctx, "UpdateItem", trace.WithAttributes(
		attribute.String("item.id", item.ID),
	))
	defer span.End()

	// Validate
	if err := item.Validate(); err != nil {
		span.RecordError(err)
		return err
	}

	item.UpdatedAt = time.Now()

	// Update
	if err := s.repo.Update(ctx, item); err != nil {
		span.RecordError(err)
		return err
	}

	// Publish event
	event := &domain.ItemEvent{
		EventType: "item.updated",
		ItemID:    item.ID,
		Timestamp: time.Now(),
		Data:      *item,
	}
	if err := s.publisher.PublishItemEvent(ctx, event); err != nil {
		span.RecordError(err)
		// Log error but don't fail the operation
	}

	return nil
}

// DeleteItem deletes an item
func (s *ItemServiceImpl) DeleteItem(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "DeleteItem", trace.WithAttributes(
		attribute.String("item.id", id),
	))
	defer span.End()

	if err := s.repo.Delete(ctx, id); err != nil {
		span.RecordError(err)
		return err
	}

	// Publish event
	event := &domain.ItemEvent{
		EventType: "item.deleted",
		ItemID:    id,
		Timestamp: time.Now(),
	}
	if err := s.publisher.PublishItemEvent(ctx, event); err != nil {
		span.RecordError(err)
		// Log error but don't fail the operation
	}

	return nil
}
