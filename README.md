# Better Auth SDK for Go

A Go SDK for [Better Auth](https://www.better-auth.com/), providing a clean and idiomatic interface for authentication and session management in your Go applications.

## Features

- 🔐 **Session Management** - Verify and retrieve session information
- 🚀 **Simple API** - Clean and intuitive Go interfaces
- ⚡ **Context Support** - Built-in support for `context.Context`
- 🛡️ **Type Safety** - Fully typed structs for users and sessions
- 🔧 **Configurable** - Flexible configuration options
- 📦 **Zero Dependencies** - Uses only Go standard library
- ⚠️ **Rich Error Handling** - Detailed error types and messages

## Installation

```bash
go get github.com/Zytera/better-auth-sdk-go
```

## Requirements

- Go 1.21 or higher

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    
    betterauth "github.com/Zytera/better-auth-sdk-go"
)

func main() {
    // Configure the client
    config := &betterauth.Config{
        BaseURL: "https://your-app.com",
    }
    
    // Create a session token from an HTTP cookie
    sessionToken := &betterauth.SessionToken{
        Cookie: &http.Cookie{
            Name:  "better-auth.session_token",
            Value: "your-session-token-here",
        },
    }
    
    // Initialize the client
    client := betterauth.NewClient(config, sessionToken)
    
    // Get session data
    ctx := context.Background()
    sessionData, err := client.Session.GetSession(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User: %s (%s)\n", sessionData.User.Name, sessionData.User.Email)
    fmt.Printf("Session ID: %s\n", sessionData.Session.ID)
}
```

### HTTP Middleware Example

```go
func authMiddleware(client *betterauth.Client) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get the session cookie
            cookie, err := r.Cookie("better-auth.session_token")
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            // Update client session token
            client.SessionToken = &betterauth.SessionToken{Cookie: cookie}
            
            // Verify session
            sessionData, err := client.Session.GetSession(r.Context())
            if err != nil {
                http.Error(w, "Invalid session", http.StatusUnauthorized)
                return
            }
            
            // Add user to context and continue
            ctx := context.WithValue(r.Context(), "user", sessionData.User)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## API Reference

### Client

#### Creating a Client

```go
client := betterauth.NewClient(config, sessionToken)
```

**Parameters:**
- `config` (*Config): Client configuration
- `sessionToken` (*SessionToken): Session token from cookie

#### Client Methods

- `SetTimeout(timeout time.Duration)` - Update the HTTP client timeout

### Configuration

```go
type Config struct {
    BaseURL    string        // Base URL of your Better Auth server (required)
    Timeout    time.Duration // HTTP client timeout (default: 30s)
    HTTPClient *http.Client  // Custom HTTP client (optional)
    Debug      bool          // Enable debug logging
}
```

**Example:**

```go
config := &betterauth.Config{
    BaseURL: "https://api.yourapp.com",
    Timeout: 60 * time.Second,
    Debug:   true,
}
```

### Session Service

#### GetSession

Retrieves the current user session and user information.

```go
sessionData, err := client.Session.GetSession(ctx)
```

**Returns:** `*SessionData`, `error`

**Example Response:**

```go
type SessionData struct {
    User    User
    Session Session
}
```

#### Verify

Verifies a session token.

```go
session, err := client.Session.Verify(ctx, "token-string")
```

**Parameters:**
- `ctx` (context.Context): Request context
- `token` (string): Session token to verify

**Returns:** `*Session`, `error`

### Types

#### User

```go
type User struct {
    ID                  string
    Email               string
    EmailVerified       bool
    Name                string
    Image               *string
    PhoneNumber         string
    PhoneNumberVerified bool
    Role                string
    Banned              bool
    BanReason           *string
    BanExpires          *time.Time
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

#### Session

```go
type Session struct {
    ID                   string
    UserID               string
    Token                string
    RefreshToken         string
    ExpiresAt            time.Time
    CreatedAt            time.Time
    UpdatedAt            time.Time
    IPAddress            string
    UserAgent            string
    ActiveOrganizationID *string
    ImpersonatedBy       *string
}
```

#### SessionToken

```go
type SessionToken struct {
    Cookie *http.Cookie
}
```

## Error Handling

The SDK provides comprehensive error handling with typed errors.

### Error Types

```go
const (
    ErrorTypeValidation   ErrorType = "validation"
    ErrorTypeUnauthorized ErrorType = "unauthorized"
    ErrorTypeNotFound     ErrorType = "not_found"
    ErrorTypeForbidden    ErrorType = "forbidden"
    ErrorTypeConflict     ErrorType = "conflict"
    ErrorTypeInternal     ErrorType = "internal"
    ErrorTypeNetwork      ErrorType = "network"
    ErrorTypeTimeout      ErrorType = "timeout"
)
```

### Error Checking

```go
sessionData, err := client.Session.GetSession(ctx)
if err != nil {
    if betterauth.IsUnauthorizedError(err) {
        // Handle unauthorized error
        log.Println("User is not authenticated")
    } else if betterauth.IsNetworkError(err) {
        // Handle network error
        log.Println("Network connection failed")
    } else {
        // Handle other errors
        log.Printf("Error: %v", err)
    }
}
```

### Error Helper Functions

- `IsError(err error) bool` - Check if error is a Better Auth error
- `IsValidationError(err error) bool`
- `IsUnauthorizedError(err error) bool`
- `IsNotFoundError(err error) bool`
- `IsForbiddenError(err error) bool`
- `IsConflictError(err error) bool`
- `IsInternalError(err error) bool`
- `IsNetworkError(err error) bool`
- `IsTimeoutError(err error) bool`

### Accessing Error Details

```go
if err != nil {
    if betterAuthErr, ok := err.(*betterauth.Error); ok {
        fmt.Printf("Error Type: %s\n", betterAuthErr.Type)
        fmt.Printf("Message: %s\n", betterAuthErr.Message)
        fmt.Printf("Status Code: %d\n", betterAuthErr.StatusCode)
        fmt.Printf("Details: %v\n", betterAuthErr.Details)
    }
}
```

## Examples

### Complete Web Server Example

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    
    betterauth "github.com/Zytera/better-auth-sdk-go"
)

func main() {
    config := &betterauth.Config{
        BaseURL: "https://your-auth-server.com",
    }
    
    mux := http.NewServeMux()
    
    // Protected route
    mux.HandleFunc("/api/profile", func(w http.ResponseWriter, r *http.Request) {
        // Get session cookie
        cookie, err := r.Cookie("better-auth.session_token")
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Create client with session token
        client := betterauth.NewClient(config, &betterauth.SessionToken{
            Cookie: cookie,
        })
        
        // Get session
        sessionData, err := client.Session.GetSession(r.Context())
        if err != nil {
            if betterauth.IsUnauthorizedError(err) {
                http.Error(w, "Invalid session", http.StatusUnauthorized)
                return
            }
            http.Error(w, "Internal error", http.StatusInternalServerError)
            return
        }
        
        // Return user profile
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(sessionData.User)
    })
    
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### Custom HTTP Client Example

```go
config := &betterauth.Config{
    BaseURL: "https://your-auth-server.com",
    HTTPClient: &http.Client{
        Timeout: 10 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
        },
    },
}

client := betterauth.NewClient(config, sessionToken)
```

## Best Practices

1. **Reuse Client Instances** - Create the client once and reuse it across requests to benefit from connection pooling
2. **Use Context** - Always pass a proper context for cancellation and timeout support
3. **Handle Errors Properly** - Use the typed error helpers to handle different error scenarios
4. **Secure Cookie Handling** - Always use secure cookies in production (HttpOnly, Secure, SameSite)
5. **Set Appropriate Timeouts** - Configure timeouts based on your application's needs
6. **Validate Configuration** - Call `config.Validate()` before creating the client if needed

## Testing

```go
func TestSessionRetrieval(t *testing.T) {
    config := &betterauth.Config{
        BaseURL: "https://test-server.com",
    }
    
    sessionToken := &betterauth.SessionToken{
        Cookie: &http.Cookie{
            Name:  "better-auth.session_token",
            Value: "test-token",
        },
    }
    
    client := betterauth.NewClient(config, sessionToken)
    
    ctx := context.Background()
    sessionData, err := client.Session.GetSession(ctx)
    
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if sessionData.User.Email == "" {
        t.Error("Expected user email to be populated")
    }
}
```

## Project Structure

```
better-auth-sdk-go/
├── client.go       # Main client implementation
├── config.go       # Configuration structures
├── errors.go       # Error types and handling
├── session.go      # Session service implementation
├── types.go        # Type definitions (User, Session, etc.)
├── validation.go   # Validation utilities
├── go.mod          # Go module definition
├── LICENSE         # GPL-3.0 license
└── README.md       # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Links

- [Better Auth Documentation](https://www.better-auth.com/docs)
- [GitHub Repository](https://github.com/Zytera/better-auth-sdk-go)

## Support

For issues, questions, or contributions, please open an issue on the GitHub repository.

## Changelog

### v1.0.0 (Current)
- Initial release
- Session verification
- Session retrieval
- User information
- Comprehensive error handling
- Type-safe API