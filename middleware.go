package betterauth

import (
	"context"
	"net/http"
	"strings"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// ContextKeyUser is the context key for storing user information
	ContextKeyUser ContextKey = "better_auth_user"
	// ContextKeySession is the context key for storing session information
	ContextKeySession ContextKey = "better_auth_session"
)

// Middleware provides HTTP middleware functions for Better Auth
type Middleware struct {
	client       *Client
	extractToken func(*http.Request) string
}

// NewMiddleware creates a new Middleware instance
func NewMiddleware(client *Client) *Middleware {
	m := &Middleware{
		client: client,
	}
	m.extractToken = m.defaultExtractToken
	return m
}

// Authenticate is a middleware that verifies the session token from the request
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		session, err := m.client.Session.Verify(r.Context(), token)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		// Get user information
		user, err := m.client.User.Get(r.Context(), session.UserID)
		if err != nil {
			http.Error(w, "Unauthorized: user not found", http.StatusUnauthorized)
			return
		}

		// Add session and user to context
		ctx := context.WithValue(r.Context(), ContextKeySession, session)
		ctx = context.WithValue(ctx, ContextKeyUser, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is a middleware that requires authentication (similar to Authenticate)
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return m.Authenticate(next)
}

// OptionalAuth is a middleware that adds user/session to context if available, but doesn't require it
func (m *Middleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token != "" {
			session, err := m.client.Session.Verify(r.Context(), token)
			if err == nil {
				user, err := m.client.User.Get(r.Context(), session.UserID)
				if err == nil {
					ctx := context.WithValue(r.Context(), ContextKeySession, session)
					ctx = context.WithValue(ctx, ContextKeyUser, user)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RequireEmailVerified is a middleware that requires the user to have a verified email
func (m *Middleware) RequireEmailVerified(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			http.Error(w, "Unauthorized: user not found", http.StatusUnauthorized)
			return
		}

		if !user.EmailVerified {
			http.Error(w, "Forbidden: email not verified", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// defaultExtractToken is the default token extraction method
// Checks in order: Authorization header (Bearer token), Cookie, Query parameter
func (m *Middleware) defaultExtractToken(r *http.Request) string {
	// Check Authorization header
	auth := r.Header.Get("Authorization")
	if auth != "" {
		// Bearer token
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
		// Direct token
		return auth
	}

	// Check Cookie
	cookie, err := r.Cookie("better_auth_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Check query parameter
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	return ""
}

// GetUserFromContext retrieves the user from the request context
func GetUserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(ContextKeyUser).(*User)
	if !ok {
		return nil
	}
	return user
}

// GetSessionFromContext retrieves the session from the request context
func GetSessionFromContext(ctx context.Context) *Session {
	session, ok := ctx.Value(ContextKeySession).(*Session)
	if !ok {
		return nil
	}
	return session
}

// AuthHandler wraps a handler function with authentication
func (m *Middleware) AuthHandler(handler func(w http.ResponseWriter, r *http.Request, user *User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		session, err := m.client.Session.Verify(r.Context(), token)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		user, err := m.client.User.Get(r.Context(), session.UserID)
		if err != nil {
			http.Error(w, "Unauthorized: user not found", http.StatusUnauthorized)
			return
		}

		// Add session and user to context
		ctx := context.WithValue(r.Context(), ContextKeySession, session)
		ctx = context.WithValue(ctx, ContextKeyUser, user)

		handler(w, r.WithContext(ctx), user)
	}
}

// WithTokenExtractor allows setting a custom token extractor function
// Note: This creates a new middleware instance with the custom extractor
func WithTokenExtractor(client *Client, extractor func(*http.Request) string) *Middleware {
	m := &Middleware{
		client: client,
	}

	// Override the extractToken method with custom logic
	originalExtract := m.extractToken
	customExtract := func(r *http.Request) string {
		token := extractor(r)
		if token != "" {
			return token
		}
		return originalExtract(r)
	}

	// Create a new middleware with custom extractor
	return &Middleware{
		client:       client,
		extractToken: customExtract,
	}
}
