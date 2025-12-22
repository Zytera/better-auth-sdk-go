package betterauth

import (
	"context"
	"fmt"
)

// AuthService handles authentication operations
type AuthService struct {
	client *Client
}

// newAuthService creates a new AuthService
func newAuthService(client *Client) *AuthService {
	return &AuthService{
		client: client,
	}
}

// SignUp registers a new user
func (s *AuthService) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	if req == nil {
		return nil, NewError(ErrorTypeValidation, "sign up request is required")
	}

	if req.Email == "" {
		return nil, NewError(ErrorTypeValidation, "email is required")
	}

	if req.Password == "" {
		return nil, NewError(ErrorTypeValidation, "password is required")
	}

	var resp SignUpResponse
	if err := s.client.doRequest(ctx, "POST", "/api/auth/sign-up", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SignIn authenticates a user with email and password
func (s *AuthService) SignIn(ctx context.Context, req *SignInRequest) (*SignInResponse, error) {
	if req == nil {
		return nil, NewError(ErrorTypeValidation, "sign in request is required")
	}

	if req.Email == "" {
		return nil, NewError(ErrorTypeValidation, "email is required")
	}

	if req.Password == "" {
		return nil, NewError(ErrorTypeValidation, "password is required")
	}

	var resp SignInResponse
	if err := s.client.doRequest(ctx, "POST", "/api/auth/sign-in", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SignOut terminates a user session
func (s *AuthService) SignOut(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return NewError(ErrorTypeValidation, "session token is required")
	}

	req := map[string]string{
		"token": sessionToken,
	}

	return s.client.doRequest(ctx, "POST", "/api/auth/sign-out", req, nil)
}

// VerifyEmail verifies a user's email address
func (s *AuthService) VerifyEmail(ctx context.Context, req *VerifyEmailRequest) error {
	if req == nil || req.Token == "" {
		return NewError(ErrorTypeValidation, "verification token is required")
	}

	return s.client.doRequest(ctx, "POST", "/api/auth/verify-email", req, nil)
}

// SendVerificationEmail sends a verification email to the user
func (s *AuthService) SendVerificationEmail(ctx context.Context, email string) error {
	if email == "" {
		return NewError(ErrorTypeValidation, "email is required")
	}

	req := map[string]string{
		"email": email,
	}

	return s.client.doRequest(ctx, "POST", "/api/auth/send-verification-email", req, nil)
}

// ResetPassword initiates a password reset flow
func (s *AuthService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	if req == nil || req.Email == "" {
		return NewError(ErrorTypeValidation, "email is required")
	}

	return s.client.doRequest(ctx, "POST", "/api/auth/reset-password", req, nil)
}

// ConfirmPasswordReset confirms a password reset with token
func (s *AuthService) ConfirmPasswordReset(ctx context.Context, req *ConfirmPasswordResetRequest) error {
	if req == nil {
		return NewError(ErrorTypeValidation, "confirm password reset request is required")
	}

	if req.Token == "" {
		return NewError(ErrorTypeValidation, "reset token is required")
	}

	if req.NewPassword == "" {
		return NewError(ErrorTypeValidation, "new password is required")
	}

	return s.client.doRequest(ctx, "POST", "/api/auth/confirm-password-reset", req, nil)
}

// ChangePassword changes the user's password
func (s *AuthService) ChangePassword(ctx context.Context, sessionToken string, req *ChangePasswordRequest) error {
	if sessionToken == "" {
		return NewError(ErrorTypeValidation, "session token is required")
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

	return s.client.doRequest(ctx, "POST", "/api/auth/change-password", req, nil)
}

// GetOAuthURL generates an OAuth URL for social login
func (s *AuthService) GetOAuthURL(ctx context.Context, provider Provider, state string) (*OAuthURLResponse, error) {
	if provider == "" {
		return nil, NewError(ErrorTypeValidation, "provider is required")
	}

	path := fmt.Sprintf("/api/auth/oauth/%s/url", provider)
	req := map[string]string{
		"state": state,
	}

	var resp OAuthURLResponse
	if err := s.client.doRequest(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// HandleOAuthCallback handles the OAuth callback and creates a session
func (s *AuthService) HandleOAuthCallback(ctx context.Context, req *OAuthCallbackRequest) (*SignInResponse, error) {
	if req == nil {
		return nil, NewError(ErrorTypeValidation, "OAuth callback request is required")
	}

	if req.Provider == "" {
		return nil, NewError(ErrorTypeValidation, "provider is required")
	}

	if req.Code == "" {
		return nil, NewError(ErrorTypeValidation, "authorization code is required")
	}

	path := fmt.Sprintf("/api/auth/oauth/%s/callback", req.Provider)

	var resp SignInResponse
	if err := s.client.doRequest(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SetupTwoFactor initiates two-factor authentication setup
func (s *AuthService) SetupTwoFactor(ctx context.Context, sessionToken string, req *TwoFactorSetupRequest) (*TwoFactorSetupResponse, error) {
	if sessionToken == "" {
		return nil, NewError(ErrorTypeValidation, "session token is required")
	}

	if req == nil || req.Method == "" {
		return nil, NewError(ErrorTypeValidation, "2FA method is required")
	}

	var resp TwoFactorSetupResponse
	if err := s.client.doRequest(ctx, "POST", "/api/auth/2fa/setup", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// VerifyTwoFactor verifies a two-factor authentication code
func (s *AuthService) VerifyTwoFactor(ctx context.Context, sessionToken string, req *TwoFactorVerifyRequest) error {
	if sessionToken == "" {
		return NewError(ErrorTypeValidation, "session token is required")
	}

	if req == nil || req.Code == "" {
		return NewError(ErrorTypeValidation, "verification code is required")
	}

	return s.client.doRequest(ctx, "POST", "/api/auth/2fa/verify", req, nil)
}

// DisableTwoFactor disables two-factor authentication
func (s *AuthService) DisableTwoFactor(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return NewError(ErrorTypeValidation, "session token is required")
	}

	return s.client.doRequest(ctx, "POST", "/api/auth/2fa/disable", nil, nil)
}
