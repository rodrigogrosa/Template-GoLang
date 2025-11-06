package ports

import (
	"context"

	"github.com/rodrigogrosa/Template-GoLang/internal/core/domain"
)

// ItemRepository defines the interface for item storage
type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) error
	GetByID(ctx context.Context, id string) (*domain.Item, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Item, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id string) error
}

// ItemService defines the interface for item business logic
type ItemService interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	GetItem(ctx context.Context, id string) (*domain.Item, error)
	ListItems(ctx context.Context, limit, offset int) ([]*domain.Item, error)
	UpdateItem(ctx context.Context, item *domain.Item) error
	DeleteItem(ctx context.Context, id string) error
}

// MessageProducer defines the interface for message publishing
type MessageProducer interface {
	PublishItemCreated(ctx context.Context, item *domain.Item) error
	PublishItemUpdated(ctx context.Context, item *domain.Item) error
	PublishItemDeleted(ctx context.Context, itemID string) error
	Close() error
}

// MessageConsumer defines the interface for message consumption
type MessageConsumer interface {
	Start(ctx context.Context) error
	Close() error
}
