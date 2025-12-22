// Package betterauth provides a Go SDK for integrating Better Auth authentication
// into your Go applications.
//
// Better Auth is a modern authentication solution that provides email/password
// authentication, OAuth/social login, session management, two-factor authentication,
// and more.
//
// # Installation
//
// Install the SDK using go get:
//
//	go get github.com/medapsis/better-auth-sdk-go
//
// # Quick Start
//
// Initialize the client with your Better Auth configuration:
//
//	import betterauth "github.com/medapsis/better-auth-sdk-go"
//
//	client := betterauth.NewClient(&betterauth.Config{
//	    BaseURL:   "https://your-app.com",
//	    APIKey:    "your-api-key",
//	    SecretKey: "your-secret-key",
//	})
//
// # Authentication
//
// Sign up a new user:
//
//	ctx := context.Background()
//	resp, err := client.Auth.SignUp(ctx, &betterauth.SignUpRequest{
//	    Email:    "user@example.com",
//	    Password: "securePassword123",
//	    Name:     "John Doe",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("User ID: %s\n", resp.User.ID)
//	fmt.Printf("Session Token: %s\n", resp.Session.Token)
//
// Sign in an existing user:
//
//	resp, err := client.Auth.SignIn(ctx, &betterauth.SignInRequest{
//	    Email:    "user@example.com",
//	    Password: "securePassword123",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	sessionToken := resp.Session.Token
//
// Sign out a user:
//
//	err := client.Auth.SignOut(ctx, sessionToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Session Management
//
// Verify a session token:
//
//	session, err := client.Session.Verify(ctx, sessionToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Session expires at: %s\n", session.ExpiresAt)
//
// Refresh a session:
//
//	newSession, err := client.Session.Refresh(ctx, refreshToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Revoke a session:
//
//	err := client.Session.Revoke(ctx, sessionToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # User Management
//
// Get user information:
//
//	user, err := client.User.Get(ctx, userID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("User: %s <%s>\n", user.Name, user.Email)
//
// Update user information:
//
//	updatedUser, err := client.User.Update(ctx, userID, &betterauth.UpdateUserRequest{
//	    Name: "Jane Doe",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Delete a user:
//
//	err := client.User.Delete(ctx, userID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # OAuth/Social Authentication
//
// Get OAuth URL for social login:
//
//	urlResp, err := client.Auth.GetOAuthURL(ctx, betterauth.ProviderGoogle, "state-token")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Redirect to: %s\n", urlResp.URL)
//
// Handle OAuth callback:
//
//	resp, err := client.Auth.HandleOAuthCallback(ctx, &betterauth.OAuthCallbackRequest{
//	    Provider: betterauth.ProviderGoogle,
//	    Code:     "auth-code-from-callback",
//	    State:    "state-token",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Two-Factor Authentication
//
// Setup 2FA for a user:
//
//	setupResp, err := client.Auth.SetupTwoFactor(ctx, sessionToken, &betterauth.TwoFactorSetupRequest{
//	    Method: "totp",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("QR Code: %s\n", setupResp.QRCode)
//
// Verify 2FA code:
//
//	err := client.Auth.VerifyTwoFactor(ctx, sessionToken, &betterauth.TwoFactorVerifyRequest{
//	    Code: "123456",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # HTTP Middleware
//
// The SDK provides HTTP middleware for easy integration with Go web applications:
//
//	middleware := betterauth.NewMiddleware(client)
//
//	// Require authentication
//	http.Handle("/api/protected", middleware.Authenticate(http.HandlerFunc(handler)))
//
//	// Optional authentication
//	http.Handle("/api/optional", middleware.OptionalAuth(http.HandlerFunc(handler)))
//
//	// Require verified email
//	http.Handle("/api/verified", middleware.RequireEmailVerified(http.HandlerFunc(handler)))
//
// Access user from context:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    user := betterauth.GetUserFromContext(r.Context())
//	    session := betterauth.GetSessionFromContext(r.Context())
//	    // Use user and session...
//	}
//
// # Error Handling
//
// The SDK provides typed errors for easy error handling:
//
//	user, err := client.Auth.SignIn(ctx, req)
//	if err != nil {
//	    if betterauth.IsUnauthorizedError(err) {
//	        // Handle invalid credentials
//	        fmt.Println("Invalid email or password")
//	    } else if betterauth.IsValidationError(err) {
//	        // Handle validation errors
//	        fmt.Println("Validation error:", err)
//	    } else if betterauth.IsNotFoundError(err) {
//	        // Handle not found
//	        fmt.Println("Resource not found")
//	    } else {
//	        // Handle other errors
//	        fmt.Println("Error:", err)
//	    }
//	}
//
// # Configuration
//
// The Config struct allows you to customize the client behavior:
//
//	config := &betterauth.Config{
//	    BaseURL:    "https://your-app.com",
//	    APIKey:     "your-api-key",
//	    SecretKey:  "your-secret-key",
//	    Timeout:    30 * time.Second,    // HTTP client timeout
//	    HTTPClient: customHTTPClient,     // Custom HTTP client (optional)
//	    Debug:      true,                 // Enable debug logging
//	}
//
// # Validation
//
// The SDK includes validation utilities:
//
//	// Validate email
//	if !betterauth.ValidateEmail(email) {
//	    fmt.Println("Invalid email")
//	}
//
//	// Validate password strength
//	if err := betterauth.ValidatePassword(password); err != nil {
//	    fmt.Println("Weak password:", err)
//	}
//
// # Thread Safety
//
// The Client and all service instances are thread-safe and can be used
// concurrently from multiple goroutines. However, individual request/response
// objects should not be shared between goroutines.
//
// # Context Support
//
// All API methods accept a context.Context parameter for cancellation and
// timeout support:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	user, err := client.User.Get(ctx, userID)
//	if err != nil {
//	    if ctx.Err() == context.DeadlineExceeded {
//	        fmt.Println("Request timeout")
//	    }
//	}
//
// # Examples
//
// For more examples, see the examples directory in the repository:
//   - examples/basic_auth - Basic authentication flows
//   - examples/session_management - Session management examples
//   - examples/middleware - HTTP middleware integration
//
// # License
//
// This SDK is licensed under the MIT License. See LICENSE file for details.
//
// # Support
//
// For issues and questions, please visit:
// https://github.com/medapsis/better-auth-sdk-go/issues
package betterauth

const (
	// Version is the current version of the SDK
	Version = "0.1.0"

	// UserAgent is the HTTP User-Agent header value
	UserAgent = "better-auth-sdk-go/" + Version
)
