package ports

import (
	"context"
	"github.com/rodrigogrosa/Template-GoLang/internal/domain"
)

// ItemService defines the business logic for items
type ItemService interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	GetItem(ctx context.Context, id string) (*domain.Item, error)
	ListItems(ctx context.Context, limit, offset int) ([]*domain.Item, error)
	UpdateItem(ctx context.Context, item *domain.Item) error
	DeleteItem(ctx context.Context, id string) error
}
