# Quick Start Guide - Better Auth SDK for Go

Get started with Better Auth SDK in minutes!

## Installation

```bash
go get github.com/medapsis/better-auth-sdk-go
```

## Initialize the Client

```go
package main

import (
    "context"
    "log"
    
    betterauth "github.com/medapsis/better-auth-sdk-go"
)

func main() {
    // Create a new client
    client := betterauth.NewClient(&betterauth.Config{
        BaseURL:   "https://your-app.com",
        APIKey:    "your-api-key",
        SecretKey: "your-secret-key",
    })
    
    ctx := context.Background()
    
    // Now you can use the client!
}
```

## Basic Operations

### 1. Sign Up a New User

```go
resp, err := client.Auth.SignUp(ctx, &betterauth.SignUpRequest{
    Email:    "user@example.com",
    Password: "SecurePassword123!",
    Name:     "John Doe",
})
if err != nil {
    log.Fatal(err)
}

log.Printf("User created: %s", resp.User.ID)
log.Printf("Session token: %s", resp.Session.Token)
```

### 2. Sign In

```go
resp, err := client.Auth.SignIn(ctx, &betterauth.SignInRequest{
    Email:    "user@example.com",
    Password: "SecurePassword123!",
})
if err != nil {
    log.Fatal(err)
}

sessionToken := resp.Session.Token
```

### 3. Verify Session

```go
session, err := client.Session.Verify(ctx, sessionToken)
if err != nil {
    log.Fatal(err)
}

log.Printf("Session valid until: %s", session.ExpiresAt)
```

### 4. Get User Profile

```go
user, err := client.User.Get(ctx, userID)
if err != nil {
    log.Fatal(err)
}

log.Printf("User: %s <%s>", user.Name, user.Email)
```

### 5. Update User

```go
updatedUser, err := client.User.Update(ctx, userID, &betterauth.UpdateUserRequest{
    Name: "Jane Doe",
})
if err != nil {
    log.Fatal(err)
}
```

### 6. Sign Out

```go
err := client.Auth.SignOut(ctx, sessionToken)
if err != nil {
    log.Fatal(err)
}
```

## HTTP Middleware

Protect your HTTP endpoints with authentication:

```go
package main

import (
    "net/http"
    
    betterauth "github.com/medapsis/better-auth-sdk-go"
)

func main() {
    client := betterauth.NewClient(&betterauth.Config{
        BaseURL: "https://your-app.com",
        APIKey:  "your-api-key",
    })
    
    middleware := betterauth.NewMiddleware(client)
    
    // Public endpoint
    http.HandleFunc("/api/public", publicHandler)
    
    // Protected endpoint
    http.Handle("/api/protected", 
        middleware.Authenticate(http.HandlerFunc(protectedHandler)))
    
    http.ListenAndServe(":8080", nil)
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Public endpoint"))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    user := betterauth.GetUserFromContext(r.Context())
    w.Write([]byte("Hello, " + user.Name))
}
```

## Error Handling

```go
user, err := client.Auth.SignIn(ctx, req)
if err != nil {
    switch {
    case betterauth.IsUnauthorizedError(err):
        log.Println("Invalid credentials")
    case betterauth.IsValidationError(err):
        log.Println("Validation error:", err)
    case betterauth.IsNotFoundError(err):
        log.Println("User not found")
    default:
        log.Println("Error:", err)
    }
    return
}
```

## OAuth/Social Login

```go
// Get OAuth URL
urlResp, err := client.Auth.GetOAuthURL(ctx, betterauth.ProviderGoogle, "state-token")
if err != nil {
    log.Fatal(err)
}

// Redirect user to urlResp.URL

// Handle callback
session, err := client.Auth.HandleOAuthCallback(ctx, &betterauth.OAuthCallbackRequest{
    Provider: betterauth.ProviderGoogle,
    Code:     "auth-code-from-callback",
    State:    "state-token",
})
```

## Password Reset

```go
// Request password reset
err := client.Auth.ResetPassword(ctx, &betterauth.ResetPasswordRequest{
    Email: "user@example.com",
})

// User receives email with token

// Confirm password reset
err = client.Auth.ConfirmPasswordReset(ctx, &betterauth.ConfirmPasswordResetRequest{
    Token:       "reset-token-from-email",
    NewPassword: "NewSecurePassword456!",
})
```

## Email Verification

```go
// Send verification email
err := client.Auth.SendVerificationEmail(ctx, "user@example.com")

// User clicks link in email

// Verify email with token
err = client.Auth.VerifyEmail(ctx, &betterauth.VerifyEmailRequest{
    Token: "verification-token",
})
```

## Session Management

```go
// List all sessions
sessions, err := client.Session.List(ctx, userID)

// Refresh session
newSession, err := client.Session.Refresh(ctx, refreshToken)

// Revoke single session
err = client.Session.Revoke(ctx, sessionToken)

// Revoke all sessions
err = client.Session.RevokeAll(ctx, userID)
```

## Two-Factor Authentication

```go
// Setup 2FA
setupResp, err := client.Auth.SetupTwoFactor(ctx, sessionToken, &betterauth.TwoFactorSetupRequest{
    Method: "totp",
})
// Show setupResp.QRCode to user

// Verify 2FA code
err = client.Auth.VerifyTwoFactor(ctx, sessionToken, &betterauth.TwoFactorVerifyRequest{
    Code: "123456",
})

// Disable 2FA
err = client.Auth.DisableTwoFactor(ctx, sessionToken)
```

## Configuration Options

```go
client := betterauth.NewClient(&betterauth.Config{
    BaseURL:    "https://your-app.com",
    APIKey:     "your-api-key",
    SecretKey:  "your-secret-key",
    Timeout:    30 * time.Second,      // Default: 30s
    HTTPClient: customHTTPClient,       // Optional
    Debug:      true,                   // Optional
})
```

## Context & Timeout

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := client.User.Get(ctx, userID)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
// Call cancel() to abort request
```

## Validation Helpers

```go
// Validate email
if !betterauth.ValidateEmail(email) {
    log.Println("Invalid email")
}

// Validate password strength
if err := betterauth.ValidatePassword(password); err != nil {
    log.Println("Weak password:", err)
}

// Mask sensitive data
maskedEmail := betterauth.MaskEmail("user@example.com") // u**r@example.com
maskedToken := betterauth.MaskToken("long-token-string") // long****ring
```

## Common Patterns

### Check if session is expired

```go
if betterauth.IsSessionExpired(session.ExpiresAt) {
    // Session expired, refresh or re-authenticate
}

timeLeft := betterauth.TimeUntilExpiry(session.ExpiresAt)
log.Printf("Session expires in: %s", betterauth.FormatDuration(timeLeft))
```

### Store user metadata

```go
resp, err := client.Auth.SignUp(ctx, &betterauth.SignUpRequest{
    Email:    "user@example.com",
    Password: "SecurePassword123!",
    Name:     "John Doe",
    Metadata: map[string]interface{}{
        "role":       "admin",
        "department": "engineering",
        "preferences": map[string]interface{}{
            "theme":    "dark",
            "language": "en",
        },
    },
})
```

### OAuth providers

Available providers:
- `betterauth.ProviderGoogle`
- `betterauth.ProviderGithub`
- `betterauth.ProviderFacebook`
- `betterauth.ProviderTwitter`
- `betterauth.ProviderApple`
- `betterauth.ProviderDiscord`
- `betterauth.ProviderMicrosoft`

## Next Steps

- **Examples**: Check out the [examples](./examples) directory
- **Full Documentation**: See [README.md](./README.md)
- **API Reference**: Run `go doc -all github.com/medapsis/better-auth-sdk-go`
- **Contributing**: Read [CONTRIBUTING.md](./CONTRIBUTING.md)

## Need Help?

- 📖 [Documentation](./README.md)
- 🐛 [Report Issues](https://github.com/medapsis/better-auth-sdk-go/issues)
- 💬 [Better Auth Docs](https://www.better-auth.com/docs)

## License

MIT License - see [LICENSE](./LICENSE)