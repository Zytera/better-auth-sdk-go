# Better Auth SDK for Go

A Go SDK for integrating [Better Auth](https://www.better-auth.com/) authentication into your Go applications.

## Features

- 🔐 Complete authentication flow support
- 🚀 Easy to use and integrate
- 📦 Zero external dependencies (uses only Go standard library)
- 🛡️ Type-safe API
- 🔄 Session management
- 👤 User management
- 🎯 Social authentication support
- ✅ Email/Password authentication

## Installation

```bash
go get github.com/yourusername/better-auth-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    betterauth "github.com/yourusername/better-auth-sdk-go"
)

func main() {
    // Initialize the client
    client := betterauth.NewClient(&betterauth.Config{
        BaseURL:    "https://your-app.com",
        APIKey:     "your-api-key",
        SecretKey:  "your-secret-key",
    })

    // Sign up a new user
    ctx := context.Background()
    user, err := client.Auth.SignUp(ctx, &betterauth.SignUpRequest{
        Email:    "user@example.com",
        Password: "securePassword123",
        Name:     "John Doe",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("User created: %s", user.ID)

    // Sign in
    session, err := client.Auth.SignIn(ctx, &betterauth.SignInRequest{
        Email:    "user@example.com",
        Password: "securePassword123",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Session token: %s", session.Token)
}
```

## Usage

### Configuration

```go
config := &betterauth.Config{
    BaseURL:    "https://your-app.com",
    APIKey:     "your-api-key",
    SecretKey:  "your-secret-key",
    Timeout:    30 * time.Second, // Optional, default is 30s
}

client := betterauth.NewClient(config)
```

### Authentication

#### Sign Up

```go
user, err := client.Auth.SignUp(ctx, &betterauth.SignUpRequest{
    Email:    "user@example.com",
    Password: "securePassword123",
    Name:     "John Doe",
})
```

#### Sign In

```go
session, err := client.Auth.SignIn(ctx, &betterauth.SignInRequest{
    Email:    "user@example.com",
    Password: "securePassword123",
})
```

#### Sign Out

```go
err := client.Auth.SignOut(ctx, sessionToken)
```

### Session Management

#### Verify Session

```go
session, err := client.Session.Verify(ctx, sessionToken)
if err != nil {
    // Session is invalid
}
```

#### Refresh Session

```go
newSession, err := client.Session.Refresh(ctx, refreshToken)
```

### User Management

#### Get User

```go
user, err := client.User.Get(ctx, userID)
```

#### Update User

```go
updatedUser, err := client.User.Update(ctx, userID, &betterauth.UpdateUserRequest{
    Name:  "Jane Doe",
    Email: "jane@example.com",
})
```

#### Delete User

```go
err := client.User.Delete(ctx, userID)
```

### Social Authentication

```go
// Get OAuth URL for social login
url, err := client.Auth.GetOAuthURL(ctx, betterauth.ProviderGoogle, "state-token")

// Handle OAuth callback
session, err := client.Auth.HandleOAuthCallback(ctx, &betterauth.OAuthCallbackRequest{
    Provider: betterauth.ProviderGoogle,
    Code:     "auth-code",
    State:    "state-token",
})
```

## Error Handling

The SDK provides structured error types:

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

## Examples

See the [examples](./examples) directory for more detailed examples:

- [Basic Authentication](./examples/basic_auth)
- [Session Management](./examples/session_management)
- [User Management](./examples/user_management)
- [Social Login](./examples/social_login)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- 📧 Email: support@example.com
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/better-auth-sdk-go/issues)
- 📖 Documentation: [Better Auth Docs](https://www.better-auth.com/docs)
