// Package apikey is the client-side plugin for the Better Auth api-key plugin:
// create, get, update, delete, list and verify API keys. All routes live under
// /api-key.
//
// To authenticate requests *with* a key (rather than manage keys), use
// client.SetAPIKey(key) on the core client instead — that sends the x-api-key
// header on every request.
package apikey

import (
	"context"
	"net/url"
	"strconv"

	betterauth "github.com/Zytera/better-auth-sdk-go"
)

const routePrefix = "/api-key"

// Plugin talks to the api-key endpoints. Construct it with New(client).
type Plugin struct {
	r betterauth.Requester
}

// New wires the plugin to any Requester (typically *betterauth.Client).
func New(r betterauth.Requester) *Plugin {
	return &Plugin{r: r}
}

func requireString(field, value string) error {
	if value == "" {
		return betterauth.NewError(betterauth.ErrorTypeValidation, field+" is required")
	}
	return nil
}

// Create issues a new API key. The returned ApiKey.Key holds the plaintext key
// and is only available here — store it now, it is never returned again.
func (p *Plugin) Create(ctx context.Context, in CreateInput) (*ApiKey, error) {
	return do[ApiKey](p.r, ctx, "POST", routePrefix+"/create", in)
}

// Get returns a key's metadata by id (the plaintext key is not included).
func (p *Plugin) Get(ctx context.Context, id string) (*ApiKey, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("id", id)
	return do[ApiKey](p.r, ctx, "GET", routePrefix+"/get?"+v.Encode(), nil)
}

// Update changes a key's mutable fields. UpdateInput.KeyID is required.
func (p *Plugin) Update(ctx context.Context, in UpdateInput) (*ApiKey, error) {
	if err := requireString("keyId", in.KeyID); err != nil {
		return nil, err
	}
	return do[ApiKey](p.r, ctx, "POST", routePrefix+"/update", in)
}

// Delete removes a key.
func (p *Plugin) Delete(ctx context.Context, keyID string) error {
	if err := requireString("keyId", keyID); err != nil {
		return err
	}
	var out struct {
		Success bool `json:"success"`
	}
	return p.r.Do(ctx, "POST", routePrefix+"/delete", map[string]string{"keyId": keyID}, &out)
}

// List returns the caller's API keys. Pass a zero-value query for defaults.
func (p *Plugin) List(ctx context.Context, q ListQuery) (*ListResult, error) {
	v := url.Values{}
	if q.Limit > 0 {
		v.Set("limit", strconv.Itoa(q.Limit))
	}
	if q.Offset > 0 {
		v.Set("offset", strconv.Itoa(q.Offset))
	}
	if q.SortBy != "" {
		v.Set("sortBy", q.SortBy)
	}
	if q.SortDirection != "" {
		v.Set("sortDirection", q.SortDirection)
	}
	path := routePrefix + "/list"
	if len(v) > 0 {
		path += "?" + v.Encode()
	}
	return do[ListResult](p.r, ctx, "GET", path, nil)
}

// Verify checks a plaintext key (optionally against required permissions).
// This is a server-side operation: it must be called from a trusted backend.
func (p *Plugin) Verify(ctx context.Context, key string, permissions map[string][]string) (*VerifyResult, error) {
	if err := requireString("key", key); err != nil {
		return nil, err
	}
	body := map[string]interface{}{"key": key}
	if len(permissions) > 0 {
		body["permissions"] = permissions
	}
	return do[VerifyResult](p.r, ctx, "POST", routePrefix+"/verify", body)
}

// DeleteAllExpired deletes every expired key. Server-side maintenance operation.
func (p *Plugin) DeleteAllExpired(ctx context.Context) error {
	return p.r.Do(ctx, "POST", routePrefix+"/delete-all-expired", map[string]interface{}{}, nil)
}

// do performs the request and decodes the JSON body into *T.
func do[T any](r betterauth.Requester, ctx context.Context, method, path string, body interface{}) (*T, error) {
	var out T
	if err := r.Do(ctx, method, path, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
