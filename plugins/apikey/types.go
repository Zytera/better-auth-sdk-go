package apikey

import "time"

// ApiKey is a single API key. The plaintext Key is only populated by Create;
// every other endpoint returns it empty (the server omits it).
type ApiKey struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Start               string                 `json:"start"`
	Prefix              string                 `json:"prefix"`
	Key                 string                 `json:"key,omitempty"` // only on Create
	UserID              string                 `json:"userId"`
	RefillInterval      *int                   `json:"refillInterval,omitempty"`
	RefillAmount        *int                   `json:"refillAmount,omitempty"`
	LastRefillAt        *time.Time             `json:"lastRefillAt,omitempty"`
	Enabled             bool                   `json:"enabled"`
	RateLimitEnabled    bool                   `json:"rateLimitEnabled"`
	RateLimitTimeWindow *int                   `json:"rateLimitTimeWindow,omitempty"`
	RateLimitMax        *int                   `json:"rateLimitMax,omitempty"`
	RequestCount        int                    `json:"requestCount"`
	Remaining           *int                   `json:"remaining,omitempty"`
	LastRequest         *time.Time             `json:"lastRequest,omitempty"`
	ExpiresAt           *time.Time             `json:"expiresAt,omitempty"`
	CreatedAt           time.Time              `json:"createdAt"`
	UpdatedAt           time.Time              `json:"updatedAt"`
	Permissions         map[string][]string    `json:"permissions,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// CreateInput is the payload for Create. All fields are optional; the
// server-only fields (Remaining, Refill*, RateLimit*, Permissions) require the
// request to be authenticated with the server secret.
type CreateInput struct {
	Name                string                 `json:"name,omitempty"`
	ExpiresIn           *int                   `json:"expiresIn,omitempty"` // seconds; nil = no expiry
	Prefix              string                 `json:"prefix,omitempty"`
	Remaining           *int                   `json:"remaining,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	RefillAmount        *int                   `json:"refillAmount,omitempty"`
	RefillInterval      *int                   `json:"refillInterval,omitempty"`
	RateLimitEnabled    *bool                  `json:"rateLimitEnabled,omitempty"`
	RateLimitTimeWindow *int                   `json:"rateLimitTimeWindow,omitempty"`
	RateLimitMax        *int                   `json:"rateLimitMax,omitempty"`
	Permissions         map[string][]string    `json:"permissions,omitempty"`
}

// UpdateInput is the payload for Update. KeyId is required; the rest are the
// mutable fields (most are server-only, like CreateInput).
type UpdateInput struct {
	KeyID               string                 `json:"keyId"`
	Name                string                 `json:"name,omitempty"`
	Enabled             *bool                  `json:"enabled,omitempty"`
	Remaining           *int                   `json:"remaining,omitempty"`
	RefillAmount        *int                   `json:"refillAmount,omitempty"`
	RefillInterval      *int                   `json:"refillInterval,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	ExpiresIn           *int                   `json:"expiresIn,omitempty"`
	RateLimitEnabled    *bool                  `json:"rateLimitEnabled,omitempty"`
	RateLimitTimeWindow *int                   `json:"rateLimitTimeWindow,omitempty"`
	RateLimitMax        *int                   `json:"rateLimitMax,omitempty"`
	Permissions         map[string][]string    `json:"permissions,omitempty"`
}

// ListQuery holds the optional pagination/sorting for List.
type ListQuery struct {
	Limit         int
	Offset        int
	SortBy        string // e.g. "createdAt"
	SortDirection string // "asc" | "desc"
}

// ListResult is the response of List.
type ListResult struct {
	APIKeys []ApiKey `json:"apiKeys"`
	Total   int      `json:"total"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
}

// VerifyError is the failure detail from Verify when Valid is false.
type VerifyError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// VerifyResult is the response of Verify.
type VerifyResult struct {
	Valid bool         `json:"valid"`
	Error *VerifyError `json:"error"`
	Key   *ApiKey      `json:"key"`
}
