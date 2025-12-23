package betterauth

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(ErrorTypeValidation, "validation failed")

	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected error type %s, got %s", ErrorTypeValidation, err.Type)
	}

	if err.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got '%s'", err.Message)
	}
}

func TestNewErrorWithDetails(t *testing.T) {
	details := map[string]interface{}{
		"field":  "email",
		"reason": "invalid format",
	}

	err := NewErrorWithDetails(ErrorTypeValidation, "validation failed", details)

	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected error type %s, got %s", ErrorTypeValidation, err.Type)
	}

	if err.Details == nil {
		t.Error("Expected details to be set")
	}

	if err.Details["field"] != "email" {
		t.Errorf("Expected field 'email', got '%v'", err.Details["field"])
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(ErrorTypeInternal, "wrapped error", originalErr)

	if wrappedErr.Type != ErrorTypeInternal {
		t.Errorf("Expected error type %s, got %s", ErrorTypeInternal, wrappedErr.Type)
	}

	if wrappedErr.Message != "wrapped error" {
		t.Errorf("Expected message 'wrapped error', got '%s'", wrappedErr.Message)
	}

	if wrappedErr.Err != originalErr {
		t.Error("Expected wrapped error to contain original error")
	}

	// Test Unwrap
	unwrapped := wrappedErr.Unwrap()
	if unwrapped != originalErr {
		t.Error("Unwrap should return the original error")
	}
}

func TestErrorString(t *testing.T) {
	err := NewError(ErrorTypeValidation, "test error")
	expected := "validation: test error"

	if err.Error() != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
	}

	// Test with wrapped error
	originalErr := errors.New("original")
	wrappedErr := WrapError(ErrorTypeInternal, "wrapped", originalErr)
	expectedWrapped := "internal: wrapped (caused by: original)"

	if wrappedErr.Error() != expectedWrapped {
		t.Errorf("Expected error string '%s', got '%s'", expectedWrapped, wrappedErr.Error())
	}
}

func TestIsError(t *testing.T) {
	err := NewError(ErrorTypeValidation, "test")

	if !IsError(err) {
		t.Error("IsError should return true for Better Auth errors")
	}

	standardErr := errors.New("standard error")
	if IsError(standardErr) {
		t.Error("IsError should return false for standard errors")
	}
}

func TestIsValidationError(t *testing.T) {
	validationErr := NewError(ErrorTypeValidation, "validation failed")

	if !IsValidationError(validationErr) {
		t.Error("IsValidationError should return true for validation errors")
	}

	otherErr := NewError(ErrorTypeUnauthorized, "unauthorized")
	if IsValidationError(otherErr) {
		t.Error("IsValidationError should return false for non-validation errors")
	}

	standardErr := errors.New("standard error")
	if IsValidationError(standardErr) {
		t.Error("IsValidationError should return false for standard errors")
	}
}

func TestIsUnauthorizedError(t *testing.T) {
	unauthorizedErr := NewError(ErrorTypeUnauthorized, "unauthorized")

	if !IsUnauthorizedError(unauthorizedErr) {
		t.Error("IsUnauthorizedError should return true for unauthorized errors")
	}

	otherErr := NewError(ErrorTypeValidation, "validation failed")
	if IsUnauthorizedError(otherErr) {
		t.Error("IsUnauthorizedError should return false for non-unauthorized errors")
	}
}

func TestIsNotFoundError(t *testing.T) {
	notFoundErr := NewError(ErrorTypeNotFound, "not found")

	if !IsNotFoundError(notFoundErr) {
		t.Error("IsNotFoundError should return true for not found errors")
	}

	otherErr := NewError(ErrorTypeValidation, "validation failed")
	if IsNotFoundError(otherErr) {
		t.Error("IsNotFoundError should return false for non-not-found errors")
	}
}

func TestIsForbiddenError(t *testing.T) {
	forbiddenErr := NewError(ErrorTypeForbidden, "forbidden")

	if !IsForbiddenError(forbiddenErr) {
		t.Error("IsForbiddenError should return true for forbidden errors")
	}

	otherErr := NewError(ErrorTypeValidation, "validation failed")
	if IsForbiddenError(otherErr) {
		t.Error("IsForbiddenError should return false for non-forbidden errors")
	}
}

func TestIsConflictError(t *testing.T) {
	conflictErr := NewError(ErrorTypeConflict, "conflict")

	if !IsConflictError(conflictErr) {
		t.Error("IsConflictError should return true for conflict errors")
	}

	otherErr := NewError(ErrorTypeValidation, "validation failed")
	if IsConflictError(otherErr) {
		t.Error("IsConflictError should return false for non-conflict errors")
	}
}

func TestIsInternalError(t *testing.T) {
	internalErr := NewError(ErrorTypeInternal, "internal error")

	if !IsInternalError(internalErr) {
		t.Error("IsInternalError should return true for internal errors")
	}

	otherErr := NewError(ErrorTypeValidation, "validation failed")
	if IsInternalError(otherErr) {
		t.Error("IsInternalError should return false for non-internal errors")
	}
}

func TestIsNetworkError(t *testing.T) {
	networkErr := NewError(ErrorTypeNetwork, "network error")

	if !IsNetworkError(networkErr) {
		t.Error("IsNetworkError should return true for network errors")
	}

	otherErr := NewError(ErrorTypeValidation, "validation failed")
	if IsNetworkError(otherErr) {
		t.Error("IsNetworkError should return false for non-network errors")
	}
}

func TestIsTimeoutError(t *testing.T) {
	timeoutErr := NewError(ErrorTypeTimeout, "timeout")

	if !IsTimeoutError(timeoutErr) {
		t.Error("IsTimeoutError should return true for timeout errors")
	}

	otherErr := NewError(ErrorTypeValidation, "validation failed")
	if IsTimeoutError(otherErr) {
		t.Error("IsTimeoutError should return false for non-timeout errors")
	}
}

func TestParseErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
		wantType   ErrorType
	}{
		{
			name:       "Bad Request",
			statusCode: 400,
			body:       []byte(`{"error":"bad_request","message":"Invalid input"}`),
			wantType:   ErrorTypeValidation,
		},
		{
			name:       "Unauthorized",
			statusCode: 401,
			body:       []byte(`{"error":"unauthorized","message":"Invalid credentials"}`),
			wantType:   ErrorTypeUnauthorized,
		},
		{
			name:       "Forbidden",
			statusCode: 403,
			body:       []byte(`{"error":"forbidden","message":"Access denied"}`),
			wantType:   ErrorTypeForbidden,
		},
		{
			name:       "Not Found",
			statusCode: 404,
			body:       []byte(`{"error":"not_found","message":"Resource not found"}`),
			wantType:   ErrorTypeNotFound,
		},
		{
			name:       "Conflict",
			statusCode: 409,
			body:       []byte(`{"error":"conflict","message":"Resource already exists"}`),
			wantType:   ErrorTypeConflict,
		},
		{
			name:       "Timeout",
			statusCode: 408,
			body:       []byte(`{"error":"timeout","message":"Request timeout"}`),
			wantType:   ErrorTypeTimeout,
		},
		{
			name:       "Internal Server Error",
			statusCode: 500,
			body:       []byte(`{"error":"internal","message":"Internal server error"}`),
			wantType:   ErrorTypeInternal,
		},
		{
			name:       "Invalid JSON",
			statusCode: 500,
			body:       []byte(`invalid json`),
			wantType:   ErrorTypeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseErrorResponse(tt.statusCode, tt.body)

			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			betterAuthErr, ok := err.(*Error)
			if !ok {
				t.Fatal("Expected Better Auth error")
			}

			if betterAuthErr.Type != tt.wantType {
				t.Errorf("Expected error type %s, got %s", tt.wantType, betterAuthErr.Type)
			}

			if betterAuthErr.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, betterAuthErr.StatusCode)
			}
		})
	}
}
