package domain

import (
	"testing"
)

func TestItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    Item
		wantErr bool
		errMsg  error
	}{
		{
			name: "valid item",
			item: Item{
				Name:  "Test Item",
				Price: 10.50,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			item: Item{
				Name:  "",
				Price: 10.50,
			},
			wantErr: true,
			errMsg:  ErrInvalidItem,
		},
		{
			name: "negative price",
			item: Item{
				Name:  "Test Item",
				Price: -5.00,
			},
			wantErr: true,
			errMsg:  ErrInvalidPrice,
		},
		{
			name: "zero price",
			item: Item{
				Name:  "Test Item",
				Price: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != tt.errMsg {
				t.Errorf("Item.Validate() error = %v, want %v", err, tt.errMsg)
			}
		})
	}
}
