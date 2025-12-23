package betterauth

import (
	"fmt"
)

// ValidationError represents a validation error with field information
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

// Error implements the error interface for ValidationError
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
