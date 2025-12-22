package betterauth

import (
	"context"
	"fmt"
)

// UserService handles user-related operations
type UserService struct {
	client *Client
}

// newUserService creates a new UserService
func newUserService(client *Client) *UserService {
	return &UserService{
		client: client,
	}
}

// Get retrieves a user by ID
func (s *UserService) Get(ctx context.Context, userID string) (*User, error) {
	if userID == "" {
		return nil, NewError(ErrorTypeValidation, "user ID is required")
	}

	var user User
	err := s.client.doRequest(ctx, "GET", fmt.Sprintf("/api/auth/user/%s", userID), nil, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, NewError(ErrorTypeValidation, "email is required")
	}

	var user User
	err := s.client.doRequest(ctx, "GET", fmt.Sprintf("/api/auth/user/email/%s", email), nil, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates a user's information
func (s *UserService) Update(ctx context.Context, userID string, req *UpdateUserRequest) (*User, error) {
	if userID == "" {
		return nil, NewError(ErrorTypeValidation, "user ID is required")
	}

	if req == nil {
		return nil, NewError(ErrorTypeValidation, "update request is required")
	}

	var user User
	err := s.client.doRequest(ctx, "PATCH", fmt.Sprintf("/api/auth/user/%s", userID), req, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return NewError(ErrorTypeValidation, "user ID is required")
	}

	return s.client.doRequest(ctx, "DELETE", fmt.Sprintf("/api/auth/user/%s", userID), nil, nil)
}

// List retrieves a list of users with optional filters
func (s *UserService) List(ctx context.Context, opts *ListUsersOptions) (*ListUsersResponse, error) {
	if opts == nil {
		opts = &ListUsersOptions{
			Limit:  50,
			Offset: 0,
		}
	}

	var response ListUsersResponse
	err := s.client.doRequest(ctx, "GET", "/api/auth/users", opts, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error {
	if userID == "" {
		return NewError(ErrorTypeValidation, "user ID is required")
	}

	if req == nil {
		return NewError(ErrorTypeValidation, "change password request is required")
	}

	if req.CurrentPassword == "" {
		return NewError(ErrorTypeValidation, "current password is required")
	}

	if req.NewPassword == "" {
		return NewError(ErrorTypeValidation, "new password is required")
	}

	return s.client.doRequest(ctx, "POST", fmt.Sprintf("/api/auth/user/%s/change-password", userID), req, nil)
}

// VerifyEmail verifies a user's email address
func (s *UserService) VerifyEmail(ctx context.Context, token string) (*User, error) {
	if token == "" {
		return nil, NewError(ErrorTypeValidation, "verification token is required")
	}

	req := &VerifyEmailRequest{
		Token: token,
	}

	var user User
	err := s.client.doRequest(ctx, "POST", "/api/auth/verify-email", req, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ResendVerificationEmail resends the email verification email
func (s *UserService) ResendVerificationEmail(ctx context.Context, email string) error {
	if email == "" {
		return NewError(ErrorTypeValidation, "email is required")
	}

	req := map[string]string{
		"email": email,
	}

	return s.client.doRequest(ctx, "POST", "/api/auth/resend-verification", req, nil)
}

// GetAccounts retrieves all linked accounts for a user
func (s *UserService) GetAccounts(ctx context.Context, userID string) ([]*Account, error) {
	if userID == "" {
		return nil, NewError(ErrorTypeValidation, "user ID is required")
	}

	var accounts []*Account
	err := s.client.doRequest(ctx, "GET", fmt.Sprintf("/api/auth/user/%s/accounts", userID), nil, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

// LinkAccount links a new OAuth account to a user
func (s *UserService) LinkAccount(ctx context.Context, userID string, provider Provider) error {
	if userID == "" {
		return NewError(ErrorTypeValidation, "user ID is required")
	}

	if provider == "" {
		return NewError(ErrorTypeValidation, "provider is required")
	}

	req := map[string]string{
		"provider": string(provider),
	}

	return s.client.doRequest(ctx, "POST", fmt.Sprintf("/api/auth/user/%s/link-account", userID), req, nil)
}

// UnlinkAccount unlinks an OAuth account from a user
func (s *UserService) UnlinkAccount(ctx context.Context, userID string, accountID string) error {
	if userID == "" {
		return NewError(ErrorTypeValidation, "user ID is required")
	}

	if accountID == "" {
		return NewError(ErrorTypeValidation, "account ID is required")
	}

	return s.client.doRequest(ctx, "DELETE", fmt.Sprintf("/api/auth/user/%s/account/%s", userID, accountID), nil, nil)
}
