package betterauth

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError represents a validation error with field information
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

// Error implements the error interface for ValidationError
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// EmailRegex is a regular expression for basic email validation
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateSignUpRequest validates a sign-up request
func ValidateSignUpRequest(req *SignUpRequest) error {
	if req == nil {
		return NewValidationError("request", "sign up request cannot be nil")
	}

	if err := ValidateEmailField(req.Email); err != nil {
		return err
	}

	if err := ValidatePasswordField(req.Password); err != nil {
		return err
	}

	if req.Name != "" {
		if err := ValidateNameField(req.Name); err != nil {
			return err
		}
	}

	return nil
}

// ValidateSignInRequest validates a sign-in request
func ValidateSignInRequest(req *SignInRequest) error {
	if req == nil {
		return NewValidationError("request", "sign in request cannot be nil")
	}

	if err := ValidateEmailField(req.Email); err != nil {
		return err
	}

	if req.Password == "" {
		return NewValidationError("password", "password is required")
	}

	return nil
}

// ValidateUpdateUserRequest validates an update user request
func ValidateUpdateUserRequest(req *UpdateUserRequest) error {
	if req == nil {
		return NewValidationError("request", "update user request cannot be nil")
	}

	if req.Email != "" {
		if err := ValidateEmailField(req.Email); err != nil {
			return err
		}
	}

	if req.Name != "" {
		if err := ValidateNameField(req.Name); err != nil {
			return err
		}
	}

	return nil
}

// ValidateEmailField validates an email field
func ValidateEmailField(email string) error {
	if email == "" {
		return NewValidationError("email", "email is required")
	}

	email = strings.TrimSpace(email)

	if len(email) > 255 {
		return NewValidationError("email", "email is too long (max 255 characters)")
	}

	if !EmailRegex.MatchString(email) {
		return NewValidationError("email", "invalid email format")
	}

	return nil
}

// ValidatePasswordField validates a password field with strength requirements
func ValidatePasswordField(password string) error {
	if password == "" {
		return NewValidationError("password", "password is required")
	}

	if len(password) < 8 {
		return NewValidationError("password", "password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return NewValidationError("password", "password is too long (max 128 characters)")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return NewValidationError("password", "password must contain at least one uppercase letter")
	}

	if !hasLower {
		return NewValidationError("password", "password must contain at least one lowercase letter")
	}

	if !hasNumber {
		return NewValidationError("password", "password must contain at least one number")
	}

	if !hasSpecial {
		return NewValidationError("password", "password must contain at least one special character")
	}

	return nil
}

// ValidatePasswordWithOptions validates a password with custom options
type PasswordOptions struct {
	MinLength      int
	MaxLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
}

// DefaultPasswordOptions returns default password validation options
func DefaultPasswordOptions() PasswordOptions {
	return PasswordOptions{
		MinLength:      8,
		MaxLength:      128,
		RequireUpper:   true,
		RequireLower:   true,
		RequireNumber:  true,
		RequireSpecial: true,
	}
}

// ValidatePasswordWithOptions validates a password with custom options
func ValidatePasswordWithOptions(password string, opts PasswordOptions) error {
	if password == "" {
		return NewValidationError("password", "password is required")
	}

	if len(password) < opts.MinLength {
		return NewValidationError("password", fmt.Sprintf("password must be at least %d characters long", opts.MinLength))
	}

	if len(password) > opts.MaxLength {
		return NewValidationError("password", fmt.Sprintf("password is too long (max %d characters)", opts.MaxLength))
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if opts.RequireUpper && !hasUpper {
		return NewValidationError("password", "password must contain at least one uppercase letter")
	}

	if opts.RequireLower && !hasLower {
		return NewValidationError("password", "password must contain at least one lowercase letter")
	}

	if opts.RequireNumber && !hasNumber {
		return NewValidationError("password", "password must contain at least one number")
	}

	if opts.RequireSpecial && !hasSpecial {
		return NewValidationError("password", "password must contain at least one special character")
	}

	return nil
}

// ValidateNameField validates a name field
func ValidateNameField(name string) error {
	if name == "" {
		return NewValidationError("name", "name is required")
	}

	name = strings.TrimSpace(name)

	if len(name) < 2 {
		return NewValidationError("name", "name must be at least 2 characters long")
	}

	if len(name) > 100 {
		return NewValidationError("name", "name is too long (max 100 characters)")
	}

	// Check for invalid characters (allow letters, spaces, hyphens, apostrophes)
	validNameRegex := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
	if !validNameRegex.MatchString(name) {
		return NewValidationError("name", "name contains invalid characters")
	}

	return nil
}

// ValidateToken validates a token string
func ValidateToken(token string) error {
	if token == "" {
		return NewValidationError("token", "token is required")
	}

	if len(token) < 10 {
		return NewValidationError("token", "token is too short")
	}

	if len(token) > 500 {
		return NewValidationError("token", "token is too long")
	}

	return nil
}

// ValidateUserID validates a user ID
func ValidateUserID(userID string) error {
	if userID == "" {
		return NewValidationError("userID", "user ID is required")
	}

	if len(userID) > 100 {
		return NewValidationError("userID", "user ID is too long")
	}

	return nil
}

// ValidateProvider validates an OAuth provider
func ValidateProvider(provider Provider) error {
	if provider == "" {
		return NewValidationError("provider", "provider is required")
	}

	if !IsValidProvider(provider) {
		return NewValidationError("provider", fmt.Sprintf("invalid provider: %s", provider))
	}

	return nil
}

// ValidateURL validates a URL string
func ValidateURL(url string) error {
	if url == "" {
		return NewValidationError("url", "URL is required")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return NewValidationError("url", "URL must start with http:// or https://")
	}

	if len(url) > 2048 {
		return NewValidationError("url", "URL is too long (max 2048 characters)")
	}

	return nil
}

// ValidateChangePasswordRequest validates a change password request
func ValidateChangePasswordRequest(req *ChangePasswordRequest) error {
	if req == nil {
		return NewValidationError("request", "change password request cannot be nil")
	}

	if req.CurrentPassword == "" {
		return NewValidationError("currentPassword", "current password is required")
	}

	if req.NewPassword == "" {
		return NewValidationError("newPassword", "new password is required")
	}

	if req.CurrentPassword == req.NewPassword {
		return NewValidationError("newPassword", "new password must be different from current password")
	}

	if err := ValidatePasswordField(req.NewPassword); err != nil {
		return err
	}

	return nil
}

// ValidateOAuthCallbackRequest validates an OAuth callback request
func ValidateOAuthCallbackRequest(req *OAuthCallbackRequest) error {
	if req == nil {
		return NewValidationError("request", "OAuth callback request cannot be nil")
	}

	if err := ValidateProvider(req.Provider); err != nil {
		return err
	}

	if req.Code == "" {
		return NewValidationError("code", "authorization code is required")
	}

	if req.State == "" {
		return NewValidationError("state", "state parameter is required")
	}

	return nil
}

// ValidateMetadata validates metadata map
func ValidateMetadata(metadata map[string]interface{}) error {
	if metadata == nil {
		return nil
	}

	// Check for reserved keys
	reservedKeys := []string{"id", "email", "password", "createdAt", "updatedAt"}
	for _, key := range reservedKeys {
		if _, exists := metadata[key]; exists {
			return NewValidationError("metadata", fmt.Sprintf("metadata cannot contain reserved key: %s", key))
		}
	}

	// Check metadata size
	if len(metadata) > 50 {
		return NewValidationError("metadata", "metadata cannot have more than 50 keys")
	}

	return nil
}
