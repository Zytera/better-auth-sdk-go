package betterauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := &Config{
		BaseURL:   "https://example.com",
		APIKey:    "test-api-key",
		SecretKey: "test-secret-key",
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.config.BaseURL != config.BaseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", config.BaseURL, client.config.BaseURL)
	}

	if client.config.APIKey != config.APIKey {
		t.Errorf("Expected APIKey to be %s, got %s", config.APIKey, client.config.APIKey)
	}

	if client.config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30s, got %v", client.config.Timeout)
	}

	if client.Auth == nil {
		t.Error("Expected Auth service to be initialized")
	}

	if client.Session == nil {
		t.Error("Expected Session service to be initialized")
	}

	if client.User == nil {
		t.Error("Expected User service to be initialized")
	}
}

func TestNewClientWithCustomTimeout(t *testing.T) {
	config := &Config{
		BaseURL: "https://example.com",
		Timeout: 60 * time.Second,
	}

	client := NewClient(config)

	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("Expected timeout to be 60s, got %v", client.httpClient.Timeout)
	}
}

func TestDoRequest_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Expected Content-Type header to be application/json")
		}

		if r.Header.Get("X-API-Key") != "test-api-key" {
			t.Error("Expected X-API-Key header to be set")
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"id":    "123",
			"email": "test@example.com",
		})
	}))
	defer server.Close()

	client := NewClient(&Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	type Response struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	var resp Response
	err := client.doRequest(context.Background(), "POST", "/api/test", map[string]string{
		"test": "data",
	}, &resp)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.ID != "123" {
		t.Errorf("Expected ID to be 123, got %s", resp.ID)
	}

	if resp.Email != "test@example.com" {
		t.Errorf("Expected email to be test@example.com, got %s", resp.Email)
	}
}

func TestDoRequest_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid credentials",
		})
	}))
	defer server.Close()

	client := NewClient(&Config{
		BaseURL: server.URL,
	})

	err := client.doRequest(context.Background(), "POST", "/api/test", nil, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !IsUnauthorizedError(err) {
		t.Errorf("Expected unauthorized error, got %v", err)
	}
}

func TestDoRequest_ValidationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid email format",
			Details: map[string]interface{}{
				"field": "email",
			},
		})
	}))
	defer server.Close()

	client := NewClient(&Config{
		BaseURL: server.URL,
	})

	err := client.doRequest(context.Background(), "POST", "/api/test", nil, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !IsValidationError(err) {
		t.Errorf("Expected validation error, got %v", err)
	}

	betterErr, ok := err.(*Error)
	if !ok {
		t.Fatal("Expected error to be *Error type")
	}

	if betterErr.Message != "Invalid email format" {
		t.Errorf("Expected message 'Invalid email format', got %s", betterErr.Message)
	}
}

func TestDoRequest_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "User not found",
		})
	}))
	defer server.Close()

	client := NewClient(&Config{
		BaseURL: server.URL,
	})

	err := client.doRequest(context.Background(), "GET", "/api/user/123", nil, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestSetTimeout(t *testing.T) {
	client := NewClient(&Config{
		BaseURL: "https://example.com",
	})

	newTimeout := 60 * time.Second
	client.SetTimeout(newTimeout)

	if client.httpClient.Timeout != newTimeout {
		t.Errorf("Expected timeout to be %v, got %v", newTimeout, client.httpClient.Timeout)
	}
}

func TestSetAPIKey(t *testing.T) {
	client := NewClient(&Config{
		BaseURL: "https://example.com",
		APIKey:  "old-key",
	})

	newKey := "new-api-key"
	client.SetAPIKey(newKey)

	if client.config.APIKey != newKey {
		t.Errorf("Expected API key to be %s, got %s", newKey, client.config.APIKey)
	}
}

func TestSetSecretKey(t *testing.T) {
	client := NewClient(&Config{
		BaseURL:   "https://example.com",
		SecretKey: "old-secret",
	})

	newSecret := "new-secret-key"
	client.SetSecretKey(newSecret)

	if client.config.SecretKey != newSecret {
		t.Errorf("Expected secret key to be %s, got %s", newSecret, client.config.SecretKey)
	}
}

func TestDoRequest_WithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(&Config{
		BaseURL: server.URL,
	})

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	err := client.doRequest(ctx, "GET", "/api/test", nil, nil)

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}
}

func TestDoRequest_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(&Config{
		BaseURL: server.URL,
	})

	err := client.doRequest(context.Background(), "DELETE", "/api/test", nil, nil)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
