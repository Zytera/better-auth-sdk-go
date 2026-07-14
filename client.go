package betterauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Requester is the minimal surface a plugin needs to talk to the server.
// Every plugin subpackage accepts one of these (the *Client satisfies it),
// which is how the SDK stays modular like better-auth's plugins.
type Requester interface {
	Do(ctx context.Context, method, path string, body, result interface{}) error
}

// Client is the main Better Auth SDK client.
type Client struct {
	config       *Config
	HTTPClient   *http.Client
	SessionToken *SessionToken
	bearerToken  string
	mu           sync.RWMutex
}

// SetBearerToken enables the bearer plugin: subsequent requests carry an
// "Authorization: Bearer <token>" header. Pass "" to disable.
// Safe for concurrent use with Do.
func (c *Client) SetBearerToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bearerToken = token
}

// SessionToken carries the session cookie from the server.
type SessionToken struct {
	Cookie *http.Cookie
}

// NewClient creates a new Better Auth client.
func NewClient(config *Config, sessionToken *SessionToken) *Client {
	config.setDefaults()

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	} else {
		// Make a shallow copy so SetTimeout does not mutate the caller's client.
		copied := *httpClient
		httpClient = &copied
	}

	return &Client{
		config:       config,
		HTTPClient:   httpClient,
		SessionToken: sessionToken,
	}
}

// Do performs an HTTP request with proper headers and error handling.
// Plugins call this via the Requester interface.
func (c *Client) Do(ctx context.Context, method, path string, body, result interface{}) error {
	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return NewError(ErrorTypeValidation, fmt.Sprintf("failed to marshal request body: %v", err))
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	fullURL, err := c.buildURL(path)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return NewError(ErrorTypeInternal, fmt.Sprintf("failed to create request: %v", err))
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set session token, if provided.
	if c.SessionToken != nil && c.SessionToken.Cookie != nil {
		cookie := c.SessionToken.Cookie
		req.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	}

	// bearer plugin: send the token as an Authorization header when set.
	c.mu.RLock()
	bearer := c.bearerToken
	c.mu.RUnlock()
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	if c.config.Debug {
		log.Printf("[betterauth] %s %s", method, fullURL)
	}

	// Perform request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return WrapError(ErrorTypeTimeout, "request timed out", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return WrapError(ErrorTypeTimeout, "request timed out", err)
		}
		return WrapError(ErrorTypeNetwork, "failed to perform request", err)
	}
	defer resp.Body.Close()

	if c.config.Debug {
		log.Printf("[betterauth] %d %s", resp.StatusCode, fullURL)
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return WrapError(ErrorTypeNetwork, "failed to read response body", err)
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		return parseErrorResponse(resp.StatusCode, respBody)
	}

	// Parse successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return NewError(ErrorTypeInternal, fmt.Sprintf("failed to unmarshal response: %v", err))
		}
	}

	return nil
}

// buildURL joins the configured BaseURL and BasePath with the provided path,
// preserving any query string that the caller attached to path.
func (c *Client) buildURL(path string) (string, error) {
	base, err := url.Parse(c.config.BaseURL)
	if err != nil {
		return "", NewError(ErrorTypeInternal, fmt.Sprintf("invalid base URL: %v", err))
	}

	rel, err := url.Parse(strings.TrimPrefix(path, "/"))
	if err != nil {
		return "", NewError(ErrorTypeInternal, fmt.Sprintf("invalid path: %v", err))
	}

	base = base.JoinPath(c.config.BasePath)
	base = base.JoinPath(rel.Path)
	base.RawQuery = rel.RawQuery
	return base.String(), nil
}

// SetTimeout updates the timeout of the client's own HTTP client copy.
// It never mutates a custom HTTPClient passed by the caller.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}
