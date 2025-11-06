package domain

import (
	"testing"
	"time"
)

func TestItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    Item
		wantErr bool
	}{
		{
			name: "valid item",
			item: Item{
				ID:          "1",
				Name:        "Test Item",
				Description: "Test Description",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid item - empty name",
			item: Item{
				ID:          "1",
				Name:        "",
				Description: "Test Description",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
