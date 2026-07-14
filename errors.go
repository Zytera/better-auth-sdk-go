package betterauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeUnauthorized represents authentication errors
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	// ErrorTypeNotFound represents resource not found errors
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeForbidden represents permission errors
	ErrorTypeForbidden ErrorType = "forbidden"
	// ErrorTypeConflict represents conflict errors (e.g., duplicate email)
	ErrorTypeConflict ErrorType = "conflict"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "internal"
	// ErrorTypeNetwork represents network-related errors
	ErrorTypeNetwork ErrorType = "network"
	// ErrorTypeTimeout represents timeout errors
	ErrorTypeTimeout ErrorType = "timeout"
)

// Error represents a Better Auth SDK error
type Error struct {
	Type       ErrorType              `json:"type"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"status_code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new Error
func NewError(errType ErrorType, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
	}
}

// NewErrorWithDetails creates a new Error with additional details
func NewErrorWithDetails(errType ErrorType, message string, details map[string]interface{}) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Details: details,
	}
}

// WrapError wraps an existing error
func WrapError(errType ErrorType, message string, err error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// IsError checks if an error is a Better Auth error
func IsError(err error) bool {
	_, ok := err.(*Error)
	return ok
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeValidation
	}
	return false
}

// IsUnauthorizedError checks if an error is an unauthorized error
func IsUnauthorizedError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeUnauthorized
	}
	return false
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeNotFound
	}
	return false
}

// IsForbiddenError checks if an error is a forbidden error
func IsForbiddenError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeForbidden
	}
	return false
}

// IsConflictError checks if an error is a conflict error
func IsConflictError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeConflict
	}
	return false
}

// IsInternalError checks if an error is an internal error
func IsInternalError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeInternal
	}
	return false
}

// IsNetworkError checks if an error is a network error
func IsNetworkError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeNetwork
	}
	return false
}

// IsTimeoutError checks if an error is a timeout error
func IsTimeoutError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeTimeout
	}
	return false
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// parseErrorResponse parses an error response from the API
func parseErrorResponse(statusCode int, body []byte) error {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		// If we can't parse the error response, return a generic error with a
		// truncated message so an HTML error page (e.g. from a proxy) does not
		// flood the caller.
		msg := strings.TrimSpace(string(body))
		if len(msg) > 512 {
			msg = msg[:512] + "..."
		}
		if msg == "" {
			msg = fmt.Sprintf("HTTP %d error", statusCode)
		}
		return &Error{
			Type:       ErrorTypeInternal,
			Message:    msg,
			StatusCode: statusCode,
		}
	}

	// Map HTTP status codes to error types
	var errType ErrorType
	switch statusCode {
	case http.StatusBadRequest:
		errType = ErrorTypeValidation
	case http.StatusUnauthorized:
		errType = ErrorTypeUnauthorized
	case http.StatusForbidden:
		errType = ErrorTypeForbidden
	case http.StatusNotFound:
		errType = ErrorTypeNotFound
	case http.StatusConflict:
		errType = ErrorTypeConflict
	case http.StatusRequestTimeout:
		errType = ErrorTypeTimeout
	default:
		errType = ErrorTypeInternal
	}

	message := errResp.Message
	if message == "" {
		message = errResp.Error
	}
	if message == "" {
		message = fmt.Sprintf("HTTP %d error", statusCode)
	}

	return &Error{
		Type:       errType,
		Message:    message,
		StatusCode: statusCode,
		Details:    errResp.Details,
	}
}
