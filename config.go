package betterauth

import (
	"net/http"
	"time"
)

// Config holds the configuration for the Better Auth client.
type Config struct {
	// BaseURL is the base URL of your Better Auth server
	BaseURL string

	// APIKey is the API key for authentication (optional)
	APIKey string

	// SecretKey is the secret key for signing requests (optional)
	SecretKey string

	// Timeout is the HTTP client timeout (default: 30 seconds)
	Timeout time.Duration

	// HTTPClient is a custom HTTP client (optional)
	HTTPClient *http.Client

	// Debug enables debug logging
	Debug bool
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return &ValidationError{
			Field:   "BaseURL",
			Message: "base URL is required",
		}
	}

	return nil
}

// setDefaults sets default values for optional fields.
func (c *Config) setDefaults() {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Timeout: c.Timeout,
		}
	}
}
