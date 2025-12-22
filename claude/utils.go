package betterauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Utility functions for the Better Auth SDK

// ValidateEmail performs basic email validation
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	// Basic email validation
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}

	if !strings.Contains(parts[1], ".") {
		return false
	}

	return true
}

// ValidatePassword checks if password meets minimum requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return NewError(ErrorTypeValidation, "password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	if !hasUpper {
		return NewError(ErrorTypeValidation, "password must contain at least one uppercase letter")
	}

	if !hasLower {
		return NewError(ErrorTypeValidation, "password must contain at least one lowercase letter")
	}

	if !hasDigit {
		return NewError(ErrorTypeValidation, "password must contain at least one digit")
	}

	return nil
}

// IsSessionExpired checks if a session has expired
func IsSessionExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// TimeUntilExpiry returns the duration until a session expires
func TimeUntilExpiry(expiresAt time.Time) time.Duration {
	return time.Until(expiresAt)
}

// GenerateState generates a random state parameter for OAuth
func GenerateState() string {
	// In production, use crypto/rand for better randomness
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
}

// SignRequest signs a request using HMAC-SHA256
func SignRequest(secretKey string, data string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// VerifySignature verifies a request signature
func VerifySignature(secretKey string, data string, signature string) bool {
	expectedSignature := SignRequest(secretKey, data)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// ParseJSONToMap parses JSON string to map
func ParseJSONToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, WrapError(ErrorTypeValidation, "failed to parse JSON", err)
	}
	return result, nil
}

// MapToJSON converts map to JSON string
func MapToJSON(data map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", WrapError(ErrorTypeValidation, "failed to convert to JSON", err)
	}
	return string(jsonBytes), nil
}

// BuildQueryString builds a URL query string from parameters
func BuildQueryString(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	return values.Encode()
}

// SanitizeEmail normalizes an email address
func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// MaskEmail masks an email for display purposes
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 2 {
		return "*@" + domain
	}

	masked := string(username[0]) + strings.Repeat("*", len(username)-2) + string(username[len(username)-1])
	return masked + "@" + domain
}

// MaskToken masks a token for logging purposes
func MaskToken(token string) string {
	if len(token) <= 8 {
		return strings.Repeat("*", len(token))
	}

	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d seconds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%d minutes", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	}
	return fmt.Sprintf("%d days", int(d.Hours()/24))
}

// MergeMetadata merges two metadata maps
func MergeMetadata(base, updates map[string]interface{}) map[string]interface{} {
	if base == nil {
		base = make(map[string]interface{})
	}

	result := make(map[string]interface{})

	// Copy base
	for k, v := range base {
		result[k] = v
	}

	// Apply updates
	for k, v := range updates {
		result[k] = v
	}

	return result
}

// IsValidProvider checks if a provider is valid
func IsValidProvider(provider Provider) bool {
	validProviders := []Provider{
		ProviderGoogle,
		ProviderGithub,
		ProviderFacebook,
		ProviderTwitter,
		ProviderApple,
		ProviderDiscord,
		ProviderMicrosoft,
	}

	for _, p := range validProviders {
		if p == provider {
			return true
		}
	}

	return false
}

// GetProviderDisplayName returns a user-friendly name for a provider
func GetProviderDisplayName(provider Provider) string {
	switch provider {
	case ProviderGoogle:
		return "Google"
	case ProviderGithub:
		return "GitHub"
	case ProviderFacebook:
		return "Facebook"
	case ProviderTwitter:
		return "Twitter"
	case ProviderApple:
		return "Apple"
	case ProviderDiscord:
		return "Discord"
	case ProviderMicrosoft:
		return "Microsoft"
	default:
		return string(provider)
	}
}

// SafeString returns a string value from a map, or empty string if not found
func SafeString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// SafeInt returns an int value from a map, or 0 if not found
func SafeInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return 0
}

// SafeBool returns a bool value from a map, or false if not found
func SafeBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
