package betterauth

import (
	"context"
	"fmt"
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

// Refresh refreshes a session using a refresh token
func (s *SessionService) Refresh(ctx context.Context, refreshToken string) (*Session, error) {
	if refreshToken == "" {
		return nil, NewError(ErrorTypeValidation, "refresh token is required")
	}

	req := RefreshSessionRequest{
		RefreshToken: refreshToken,
	}

	var session Session
	err := s.client.doRequest(ctx, "POST", "/api/auth/session/refresh", req, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// Revoke revokes a session (signs out)
func (s *SessionService) Revoke(ctx context.Context, token string) error {
	if token == "" {
		return NewError(ErrorTypeValidation, "session token is required")
	}

	err := s.client.doRequest(ctx, "POST", "/api/auth/session/revoke", map[string]string{
		"token": token,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

// RevokeAll revokes all sessions for a user
func (s *SessionService) RevokeAll(ctx context.Context, userID string) error {
	if userID == "" {
		return NewError(ErrorTypeValidation, "user ID is required")
	}

	path := fmt.Sprintf("/api/auth/session/revoke-all/%s", userID)
	err := s.client.doRequest(ctx, "POST", path, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// List returns all active sessions for a user
func (s *SessionService) List(ctx context.Context, userID string) (*ListSessionsResponse, error) {
	if userID == "" {
		return nil, NewError(ErrorTypeValidation, "user ID is required")
	}

	path := fmt.Sprintf("/api/auth/session/list/%s", userID)
	var response ListSessionsResponse
	err := s.client.doRequest(ctx, "GET", path, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCurrent retrieves the current session based on the provided token
func (s *SessionService) GetCurrent(ctx context.Context, token string) (*Session, error) {
	return s.Verify(ctx, token)
}

// Update updates session metadata
func (s *SessionService) Update(ctx context.Context, sessionID string, updates map[string]interface{}) (*Session, error) {
	if sessionID == "" {
		return nil, NewError(ErrorTypeValidation, "session ID is required")
	}

	path := fmt.Sprintf("/api/auth/session/%s", sessionID)
	var session Session
	err := s.client.doRequest(ctx, "PATCH", path, updates, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
