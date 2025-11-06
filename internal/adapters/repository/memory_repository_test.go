package repository

import (
	"context"
	"testing"

	"github.com/rodrigogrosa/Template-GoLang/internal/core/domain"
)

func TestInMemoryItemRepository_Create(t *testing.T) {
	repo := NewInMemoryItemRepository()

	item := &domain.Item{
		ID:          "1",
		Name:        "Test Item",
		Description: "Test Description",
	}

	err := repo.Create(context.Background(), item)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Try to create duplicate
	err = repo.Create(context.Background(), item)
	if err == nil {
		t.Error("Create() should return error for duplicate item")
	}
}

func TestInMemoryItemRepository_GetByID(t *testing.T) {
	repo := NewInMemoryItemRepository()

	item := &domain.Item{
		ID:          "1",
		Name:        "Test Item",
		Description: "Test Description",
	}

	_ = repo.Create(context.Background(), item)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "get existing item",
			id:      "1",
			wantErr: false,
		},
		{
			name:    "get non-existing item",
			id:      "999",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Error("GetByID() should return item")
			}
		})
	}
}

func TestInMemoryItemRepository_List(t *testing.T) {
	repo := NewInMemoryItemRepository()

	// Create test items
	for i := 0; i < 5; i++ {
		item := &domain.Item{
			ID:          string(rune(i + 1)),
			Name:        "Test Item",
			Description: "Test Description",
		}
		_ = repo.Create(context.Background(), item)
	}

	items, err := repo.List(context.Background(), 10, 0)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if len(items) == 0 {
		t.Error("List() should return items")
	}
}

func TestInMemoryItemRepository_Update(t *testing.T) {
	repo := NewInMemoryItemRepository()

	item := &domain.Item{
		ID:          "1",
		Name:        "Test Item",
		Description: "Test Description",
	}

	_ = repo.Create(context.Background(), item)

	item.Name = "Updated Item"
	err := repo.Update(context.Background(), item)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// Try to update non-existing item
	nonExisting := &domain.Item{
		ID:   "999",
		Name: "Non-existing",
	}
	err = repo.Update(context.Background(), nonExisting)
	if err == nil {
		t.Error("Update() should return error for non-existing item")
	}
}

func TestInMemoryItemRepository_Delete(t *testing.T) {
	repo := NewInMemoryItemRepository()

	item := &domain.Item{
		ID:          "1",
		Name:        "Test Item",
		Description: "Test Description",
	}

	_ = repo.Create(context.Background(), item)

	err := repo.Delete(context.Background(), "1")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(context.Background(), "1")
	if err == nil {
		t.Error("GetByID() should return error after deletion")
	}

	// Try to delete non-existing item
	err = repo.Delete(context.Background(), "999")
	if err == nil {
		t.Error("Delete() should return error for non-existing item")
	}
}
