// Package session is the client-side plugin for Better Auth session endpoints.
package session

import (
	"context"

	betterauth "github.com/Zytera/better-auth-sdk-go"
)

// Plugin talks to the session endpoints. Construct it with New(client).
type Plugin struct {
	r betterauth.Requester
}

// New wires the plugin to any Requester (typically *betterauth.Client).
func New(r betterauth.Requester) *Plugin {
	return &Plugin{r: r}
}

// Verify verifies a session token and returns the session details.
func (p *Plugin) Verify(ctx context.Context, token string) (*betterauth.Session, error) {
	if token == "" {
		return nil, betterauth.NewError(betterauth.ErrorTypeValidation, "session token is required")
	}

	var session betterauth.Session
	err := p.r.Do(ctx, "POST", "/session/verify", map[string]string{
		"token": token,
	}, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Get returns the session for the client's current token.
func (p *Plugin) Get(ctx context.Context) (*betterauth.SessionData, error) {
	var data betterauth.SessionData
	if err := p.r.Do(ctx, "GET", "/get-session", nil, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
