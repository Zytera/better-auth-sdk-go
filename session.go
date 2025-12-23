package betterauth

import (
	"context"
)

// SessionService handles session-related operations
type SessionService struct {
	client *Client
}

// newSessionService creates a new SessionService
func newSessionService(client *Client) *SessionService {
	return &SessionService{
		client: client,
	}
}

// Verify verifies a session token and returns the session details
func (s *SessionService) Verify(ctx context.Context, token string) (*Session, error) {
	if token == "" {
		return nil, NewError(ErrorTypeValidation, "session token is required")
	}

	var session Session
	err := s.client.doRequest(ctx, "POST", "/api/auth/session/verify", map[string]string{
		"token": token,
	}, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
