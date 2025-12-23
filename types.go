package betterauth

import "time"

// User represents a user in the system
type User struct {
	ID                  string     `json:"id"`
	Email               string     `json:"email"`
	EmailVerified       bool       `json:"emailVerified"`
	Name                string     `json:"name"`
	Image               *string    `json:"image"`
	PhoneNumber         string     `json:"phoneNumber"`
	PhoneNumberVerified bool       `json:"phoneNumberVerified"`
	Role                string     `json:"role"`
	Banned              bool       `json:"banned"`
	BanReason           *string    `json:"banReason"`
	BanExpires          *time.Time `json:"banExpires"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

// Session represents a user session
type Session struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"userId"`
	Token                string    `json:"token"`
	RefreshToken         string    `json:"refreshToken,omitempty"`
	ExpiresAt            time.Time `json:"expiresAt"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	IPAddress            string    `json:"ipAddress,omitempty"`
	UserAgent            string    `json:"userAgent,omitempty"`
	ActiveOrganizationID *string   `json:"activeOrganizationId"`
	ImpersonatedBy       *string   `json:"impersonatedBy"`
}
