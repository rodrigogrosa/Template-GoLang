package services

import (
	"context"
	"testing"

	"github.com/rodrigogrosa/Template-GoLang/internal/adapters/repository"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/domain"
)

func TestItemService_CreateItem(t *testing.T) {
	repo := repository.NewInMemoryItemRepository()
	service := NewItemService(repo, nil)

	tests := []struct {
		name    string
		item    *domain.Item
		wantErr bool
	}{
		{
			name: "create valid item",
			item: &domain.Item{
				Name:        "Test Item",
				Description: "Test Description",
			},
			wantErr: false,
		},
		{
			name: "create invalid item - empty name",
			item: &domain.Item{
				Name:        "",
				Description: "Test Description",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateItem(context.Background(), tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.item.ID == "" {
				t.Error("CreateItem() should set item ID")
			}
		})
	}
}

func TestItemService_GetItem(t *testing.T) {
	repo := repository.NewInMemoryItemRepository()
	service := NewItemService(repo, nil)

	// Create test item
	item := &domain.Item{
		Name:        "Test Item",
		Description: "Test Description",
	}
	err := service.CreateItem(context.Background(), item)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "get existing item",
			id:      item.ID,
			wantErr: false,
		},
		{
			name:    "get non-existing item",
			id:      "non-existent",
			wantErr: true,
		},
		{
			name:    "get with empty id",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetItem(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Error("GetItem() should return item")
			}
		})
	}
}

func TestItemService_ListItems(t *testing.T) {
	repo := repository.NewInMemoryItemRepository()
	service := NewItemService(repo, nil)

	// Create test items
	for i := 0; i < 5; i++ {
		item := &domain.Item{
			Name:        "Test Item",
			Description: "Test Description",
		}
		_ = service.CreateItem(context.Background(), item)
	}

	tests := []struct {
		name   string
		limit  int
		offset int
	}{
		{
			name:   "list with default limit",
			limit:  10,
			offset: 0,
		},
		{
			name:   "list with custom limit",
			limit:  2,
			offset: 0,
		},
		{
			name:   "list with offset",
			limit:  10,
			offset: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := service.ListItems(context.Background(), tt.limit, tt.offset)
			if err != nil {
				t.Errorf("ListItems() error = %v", err)
			}
			if len(items) > tt.limit {
				t.Errorf("ListItems() returned %d items, want max %d", len(items), tt.limit)
			}
		})
	}
}

func TestItemService_UpdateItem(t *testing.T) {
	repo := repository.NewInMemoryItemRepository()
	service := NewItemService(repo, nil)

	// Create test item
	item := &domain.Item{
		Name:        "Test Item",
		Description: "Test Description",
	}
	err := service.CreateItem(context.Background(), item)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	tests := []struct {
		name    string
		item    *domain.Item
		wantErr bool
	}{
		{
			name: "update existing item",
			item: &domain.Item{
				ID:          item.ID,
				Name:        "Updated Item",
				Description: "Updated Description",
			},
			wantErr: false,
		},
		{
			name: "update non-existing item",
			item: &domain.Item{
				ID:          "non-existent",
				Name:        "Updated Item",
				Description: "Updated Description",
			},
			wantErr: true,
		},
		{
			name: "update with invalid data",
			item: &domain.Item{
				ID:          item.ID,
				Name:        "",
				Description: "Updated Description",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateItem(context.Background(), tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestItemService_DeleteItem(t *testing.T) {
	repo := repository.NewInMemoryItemRepository()
	service := NewItemService(repo, nil)

	// Create test item
	item := &domain.Item{
		Name:        "Test Item",
		Description: "Test Description",
	}
	err := service.CreateItem(context.Background(), item)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "delete existing item",
			id:      item.ID,
			wantErr: false,
		},
		{
			name:    "delete non-existing item",
			id:      "non-existent",
			wantErr: true,
		},
		{
			name:    "delete with empty id",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteItem(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
