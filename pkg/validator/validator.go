package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validator provides input validation functions
type Validator struct {
	errors []ValidationError
}

// New creates a new validator
func New() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// Required checks if a string field is not empty
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "is required",
		})
	}
	return v
}

// MinLength checks minimum string length
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at least %d characters", min),
		})
	}
	return v
}

// MaxLength checks maximum string length
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at most %d characters", max),
		})
	}
	return v
}

// Min checks minimum numeric value
func (v *Validator) Min(field string, value, min float64) *Validator {
	if value < min {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at least %v", min),
		})
	}
	return v
}

// Max checks maximum numeric value
func (v *Validator) Max(field string, value, max float64) *Validator {
	if value > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at most %v", max),
		})
	}
	return v
}

// Email validates email format
func (v *Validator) Email(field, value string) *Validator {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "must be a valid email address",
		})
	}
	return v
}

// Pattern validates against a regex pattern
func (v *Validator) Pattern(field, value, pattern string) *Validator {
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(value) {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "format is invalid",
		})
	}
	return v
}

// Valid checks if there are any validation errors
func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

// Errors returns all validation errors
func (v *Validator) Errors() []ValidationError {
	return v.errors
}

// Error returns the first error message
func (v *Validator) Error() string {
	if len(v.errors) > 0 {
		return v.errors[0].Error()
	}
	return ""
}
