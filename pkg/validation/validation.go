package validation

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Init initializes the validator
func Init() {
	validate = validator.New()
}

// ValidateStruct validates a struct
func ValidateStruct(s interface{}) error {
	if validate == nil {
		Init()
	}
	return validate.Struct(s)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}
