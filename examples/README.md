# Better Auth SDK - Examples

This directory contains example applications demonstrating how to use the Better Auth SDK for Go.

## Available Examples

### 1. Basic Authentication (`basic_auth/`)

Demonstrates fundamental authentication operations:
- User sign up
- User sign in
- Session verification
- User profile retrieval
- User profile updates
- Password changes
- Sign out

**Run:**
```bash
cd basic_auth
go run main.go
```

### 2. Session Management (`session_management/`)

Shows how to manage user sessions:
- Verify session tokens
- Refresh sessions
- List all user sessions
- Revoke specific sessions
- Revoke all user sessions
- Update session metadata

**Run:**
```bash
cd session_management
go run main.go
```

### 3. HTTP Middleware (`middleware/`)

Demonstrates integration with Go HTTP servers:
- Public endpoints
- Protected endpoints (require authentication)
- Optional authentication endpoints
- Email verification requirements
- Custom authorization logic
- Context-based user/session retrieval

**Run:**
```bash
cd middleware
go run main.go
```

Then test with curl:
```bash
# Public endpoint
curl http://localhost:8080/api/public

# Protected endpoint (requires auth token)
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/protected

# Optional auth endpoint
curl http://localhost:8080/api/optional
```

### 4. Complete Example (`complete/`)

A comprehensive example showcasing all SDK features:
- User registration and authentication
- Session management
- User profile operations
- Password reset flow
- OAuth/Social authentication URLs
- Email verification
- Validation utilities
- Error handling
- Metadata operations

**Run:**
```bash
cd complete
go run main.go
```

## Prerequisites

Before running the examples, you need:

1. A Better Auth server running
2. Valid API credentials (API Key and Secret Key)
3. Go 1.21 or higher installed

## Configuration

Update the client configuration in each example with your actual credentials:

```go
client := betterauth.NewClient(&betterauth.Config{
    BaseURL:   "https://your-app.com",        // Your Better Auth server URL
    APIKey:    "your-api-key",                // Your API key
    SecretKey: "your-secret-key",             // Your secret key
    Timeout:   30 * time.Second,
})
```

## Common Patterns

### Error Handling

All examples demonstrate proper error handling:

```go
user, err := client.Auth.SignIn(ctx, req)
if err != nil {
    if betterauth.IsUnauthorizedError(err) {
        // Handle invalid credentials
    } else if betterauth.IsValidationError(err) {
        // Handle validation errors
    } else {
        // Handle other errors
    }
}
```

### Context Usage

All API calls accept a context for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := client.User.Get(ctx, userID)
```

### Session Token Management

Store and use session tokens securely:

```go
// After sign in
resp, err := client.Auth.SignIn(ctx, req)
sessionToken := resp.Session.Token

// Use for authenticated requests
session, err := client.Session.Verify(ctx, sessionToken)
```

## Testing the Examples

### 1. Start with Basic Auth

```bash
cd basic_auth
go run main.go
```

This will walk you through the fundamental authentication flow.

### 2. Explore Session Management

```bash
cd session_management
go run main.go
```

Learn how to manage multiple sessions and refresh tokens.

### 3. Try the Middleware

```bash
cd middleware
go run main.go
```

In another terminal:
```bash
# Test public endpoint
curl http://localhost:8080/api/public

# Sign in to get a token (you'll need to implement this or use your app)
TOKEN="your-session-token"

# Test protected endpoint
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/protected

# Test optional auth
curl http://localhost:8080/api/optional
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/optional
```

### 4. Run the Complete Example

```bash
cd complete
go run main.go
```

This demonstrates all SDK features in one comprehensive example.

## Building the Examples

To build all examples:

```bash
# From the project root
make run-examples

# Or manually
go build -o bin/basic_auth examples/basic_auth/main.go
go build -o bin/session_management examples/session_management/main.go
go build -o bin/middleware examples/middleware/main.go
go build -o bin/complete examples/complete/main.go
```

## Learn More

- **SDK Documentation**: See the main [README](../README.md)
- **API Reference**: Run `go doc -all github.com/medapsis/better-auth-sdk-go`
- **Better Auth Docs**: https://www.better-auth.com/docs
- **GitHub**: https://github.com/medapsis/better-auth-sdk-go

## Troubleshooting

### Connection Errors

If you get connection errors:
1. Verify your Better Auth server is running
2. Check the `BaseURL` in the configuration
3. Ensure your network allows connections to the server

### Authentication Errors

If you get unauthorized errors:
1. Verify your API key and secret key are correct
2. Check that the user exists in your Better Auth database
3. Ensure the session token hasn't expired

### Validation Errors

If you get validation errors:
1. Check the password meets strength requirements (min 8 chars, uppercase, lowercase, number, special char)
2. Verify email format is valid
3. Ensure required fields are not empty

## Support

For issues or questions:
- Open an issue: https://github.com/medapsis/better-auth-sdk-go/issues
- Check existing issues for solutions
- Read the documentation

## License

These examples are part of the Better Auth SDK for Go and are licensed under the MIT License.