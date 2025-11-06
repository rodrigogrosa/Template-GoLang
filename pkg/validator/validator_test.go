package validator

import (
	"testing"
)

func TestValidator_Required(t *testing.T) {
	v := New()
	v.Required("name", "")

	if v.Valid() {
		t.Error("Expected validator to be invalid for empty string")
	}

	if len(v.Errors()) != 1 {
		t.Errorf("Expected 1 error, got %d", len(v.Errors()))
	}
}

func TestValidator_MinLength(t *testing.T) {
	v := New()
	v.MinLength("name", "ab", 3)

	if v.Valid() {
		t.Error("Expected validator to be invalid for short string")
	}
}

func TestValidator_MaxLength(t *testing.T) {
	v := New()
	v.MaxLength("name", "abcdefghij", 5)

	if v.Valid() {
		t.Error("Expected validator to be invalid for long string")
	}
}

func TestValidator_Min(t *testing.T) {
	v := New()
	v.Min("price", 5.0, 10.0)

	if v.Valid() {
		t.Error("Expected validator to be invalid for small value")
	}
}

func TestValidator_Max(t *testing.T) {
	v := New()
	v.Max("price", 15.0, 10.0)

	if v.Valid() {
		t.Error("Expected validator to be invalid for large value")
	}
}

func TestValidator_Email(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"valid email", "test@example.com", true},
		{"invalid email", "invalid", false},
		{"missing @", "testexample.com", false},
		{"missing domain", "test@", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Email("email", tt.email)

			if v.Valid() != tt.valid {
				t.Errorf("Email validation for %s: got %v, want %v", tt.email, v.Valid(), tt.valid)
			}
		})
	}
}

func TestValidator_Multiple(t *testing.T) {
	v := New()
	v.Required("name", "").
		MinLength("description", "ab", 5)

	if v.Valid() {
		t.Error("Expected validator to be invalid")
	}

	if len(v.Errors()) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(v.Errors()))
	}
}
