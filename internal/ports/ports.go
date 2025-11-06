package ports

import (
	"context"

	"github.com/rodrigogrosa/Template-GoLang/internal/domain"
)

// ItemRepository defines the interface for item persistence
type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) error
	GetByID(ctx context.Context, id string) (*domain.Item, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Item, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id string) error
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	PublishItemEvent(ctx context.Context, event *domain.ItemEvent) error
	Close() error
}

// EventConsumer defines the interface for consuming domain events
type EventConsumer interface {
	Subscribe(ctx context.Context, handler func(*domain.ItemEvent) error) error
	Close() error
}
