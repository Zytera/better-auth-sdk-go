package betterauth

import "time"

// User represents a Better Auth user
type User struct {
	ID            string                 `json:"id"`
	Email         string                 `json:"email"`
	Name          string                 `json:"name"`
	Image         string                 `json:"image,omitempty"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	EmailVerified bool                   `json:"emailVerified"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Session represents an authentication session
type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	IPAddress    string    `json:"ipAddress,omitempty"`
	UserAgent    string    `json:"userAgent,omitempty"`
}

// SignUpRequest represents a sign-up request
type SignUpRequest struct {
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Name     string                 `json:"name,omitempty"`
	Image    string                 `json:"image,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SignUpResponse represents a sign-up response
type SignUpResponse struct {
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

// SignInRequest represents a sign-in request
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignInResponse represents a sign-in response
type SignInResponse struct {
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	Name     string                 `json:"name,omitempty"`
	Email    string                 `json:"email,omitempty"`
	Image    string                 `json:"image,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// VerifyEmailRequest represents an email verification request
type VerifyEmailRequest struct {
	Token string `json:"token"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Email string `json:"email"`
}

// ConfirmPasswordResetRequest represents a password reset confirmation
type ConfirmPasswordResetRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// Provider represents OAuth providers
type Provider string

const (
	ProviderGoogle    Provider = "google"
	ProviderGithub    Provider = "github"
	ProviderFacebook  Provider = "facebook"
	ProviderTwitter   Provider = "twitter"
	ProviderApple     Provider = "apple"
	ProviderDiscord   Provider = "discord"
	ProviderMicrosoft Provider = "microsoft"
)

// OAuthCallbackRequest represents an OAuth callback request
type OAuthCallbackRequest struct {
	Provider Provider `json:"provider"`
	Code     string   `json:"code"`
	State    string   `json:"state"`
}

// OAuthURLResponse represents an OAuth URL response
type OAuthURLResponse struct {
	URL   string `json:"url"`
	State string `json:"state"`
}

// RefreshSessionRequest represents a session refresh request
type RefreshSessionRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// Account represents a linked account (OAuth)
type Account struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	Provider     Provider  `json:"provider"`
	ProviderID   string    `json:"providerId"`
	AccessToken  string    `json:"accessToken,omitempty"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// TwoFactorSetupRequest represents a 2FA setup request
type TwoFactorSetupRequest struct {
	Method string `json:"method"` // "totp" or "sms"
}

// TwoFactorSetupResponse represents a 2FA setup response
type TwoFactorSetupResponse struct {
	Secret string `json:"secret,omitempty"`
	QRCode string `json:"qrCode,omitempty"`
}

// TwoFactorVerifyRequest represents a 2FA verification request
type TwoFactorVerifyRequest struct {
	Code string `json:"code"`
}

// ListUsersOptions represents options for listing users
type ListUsersOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Search string `json:"search,omitempty"`
}

// ListUsersResponse represents a paginated list of users
type ListUsersResponse struct {
	Users   []*User `json:"users"`
	Total   int     `json:"total"`
	Limit   int     `json:"limit"`
	Offset  int     `json:"offset"`
	HasMore bool    `json:"hasMore"`
}

// ListSessionsResponse represents a list of sessions for a user
type ListSessionsResponse struct {
	Sessions []*Session `json:"sessions"`
	Total    int        `json:"total"`
}
