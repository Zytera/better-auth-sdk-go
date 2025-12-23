package betterauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the main Better Auth SDK client
type Client struct {
	config       *Config
	HTTPClient   *http.Client
	SessionToken *SessionToken

	// Service modules
	Session *SessionService
}

type SessionToken struct {
	Cookie *http.Cookie
}

// NewClient creates a new Better Auth client
func NewClient(config *Config, sessionToken *SessionToken) *Client {
	config.setDefaults()

	client := &Client{
		config:       config,
		HTTPClient:   config.HTTPClient,
		SessionToken: sessionToken,
	}

	// Initialize services
	client.Session = newSessionService(client)

	return client
}

// doRequest performs an HTTP request with proper headers and error handling
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.config.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set session token
	cookie := c.SessionToken.Cookie
	if cookie != nil {
		req.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	}

	// Perform request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		return parseErrorResponse(resp.StatusCode, respBody)
	}

	// Parse successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// SetTimeout updates the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}
