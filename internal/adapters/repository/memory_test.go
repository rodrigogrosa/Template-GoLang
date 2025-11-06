package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/rodrigogrosa/Template-GoLang/internal/domain"
)

func TestInMemoryRepository_Create(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	item := &domain.Item{
		ID:    "test-id",
		Name:  "Test Item",
		Price: 10.50,
	}

	err := repo.Create(ctx, item)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Verify item was created
	retrieved, err := repo.GetByID(ctx, "test-id")
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}

	if retrieved.Name != item.Name {
		t.Errorf("GetByID() name = %v, want %v", retrieved.Name, item.Name)
	}
}

func TestInMemoryRepository_GetByID_NotFound(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent")
	if err != domain.ErrItemNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, domain.ErrItemNotFound)
	}
}

func TestInMemoryRepository_Update(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	item := &domain.Item{
		ID:    "test-id",
		Name:  "Test Item",
		Price: 10.50,
	}

	_ = repo.Create(ctx, item)

	// Update item
	item.Name = "Updated Item"
	err := repo.Update(ctx, item)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// Verify update
	retrieved, _ := repo.GetByID(ctx, "test-id")
	if retrieved.Name != "Updated Item" {
		t.Errorf("Update() name = %v, want %v", retrieved.Name, "Updated Item")
	}
}

func TestInMemoryRepository_Delete(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	item := &domain.Item{
		ID:    "test-id",
		Name:  "Test Item",
		Price: 10.50,
	}

	_ = repo.Create(ctx, item)

	// Delete item
	err := repo.Delete(ctx, "test-id")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, "test-id")
	if err != domain.ErrItemNotFound {
		t.Errorf("GetByID() after delete error = %v, want %v", err, domain.ErrItemNotFound)
	}
}

func TestInMemoryRepository_List(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	// Create multiple items
	for i := 0; i < 5; i++ {
		item := &domain.Item{
			ID:    fmt.Sprintf("item-%d", i),
			Name:  "Test Item",
			Price: 10.50,
		}
		_ = repo.Create(ctx, item)
	}

	// List items
	items, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	if len(items) != 5 {
		t.Errorf("List() count = %v, want %v", len(items), 5)
	}
}
