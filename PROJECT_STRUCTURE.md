# Better Auth SDK for Go - Project Structure

This document provides an overview of the project structure and organization.

## Project Overview

Better Auth SDK for Go is a complete client library for integrating Better Auth authentication into Go applications. The SDK provides:

- Email/Password authentication
- OAuth/Social login support
- Session management
- User management
- Two-factor authentication
- HTTP middleware for easy integration
- Comprehensive error handling
- Extensive validation utilities

## Directory Structure

```
better-auth-sdk-go/
├── .github/
│   └── workflows/
│       └── test.yml           # GitHub Actions CI/CD configuration
├── examples/                  # Example applications
│   ├── basic_auth/           # Basic authentication flows
│   ├── complete/             # Comprehensive example
│   ├── middleware/           # HTTP middleware integration
│   ├── session_management/   # Session handling examples
│   └── README.md             # Examples documentation
├── auth.go                   # Authentication service (227 lines)
├── client.go                 # Main SDK client (111 lines)
├── client_test.go            # Client tests (282 lines)
├── config.go                 # Configuration (52 lines)
├── doc.go                    # Package documentation (267 lines)
├── errors.go                 # Error types and handling (201 lines)
├── errors_test.go            # Error handling tests (283 lines)
├── middleware.go             # HTTP middleware (206 lines)
├── session.go                # Session management service (122 lines)
├── types.go                  # Data types and structures (167 lines)
├── user.go                   # User management service (192 lines)
├── utils.go                  # Utility functions (275 lines)
├── validation.go             # Validation utilities (392 lines)
├── go.mod                    # Go module definition
├── Makefile                  # Build and test automation
├── README.md                 # Main documentation (201 lines)
├── QUICKSTART.md             # Quick start guide (357 lines)
├── CHANGELOG.md              # Version history (95 lines)
├── CONTRIBUTING.md           # Contributing guidelines (307 lines)
├── LICENSE                   # MIT License
└── .golangci.yml             # Linter configuration

Total Lines of Code:
- Core SDK: ~2,779 lines
- Tests: ~565 lines
- Examples: ~874 lines
- Documentation: ~959 lines
```

## Core Components

### 1. Client (`client.go`)

The main entry point for the SDK. Provides:
- HTTP request handling
- Error parsing
- Configuration management
- Service initialization

**Key Functions:**
- `NewClient()` - Create a new SDK client
- `doRequest()` - Internal HTTP request handler
- `SetTimeout()`, `SetAPIKey()`, `SetSecretKey()` - Configuration updates

### 2. Authentication Service (`auth.go`)

Handles all authentication operations:
- `SignUp()` - User registration
- `SignIn()` - User login
- `SignOut()` - User logout
- `VerifyEmail()` - Email verification
- `ResetPassword()` - Password reset initiation
- `ConfirmPasswordReset()` - Password reset confirmation
- `ChangePassword()` - Password change
- `GetOAuthURL()` - OAuth URL generation
- `HandleOAuthCallback()` - OAuth callback handling
- `SetupTwoFactor()` - 2FA setup
- `VerifyTwoFactor()` - 2FA verification
- `DisableTwoFactor()` - 2FA disable

### 3. Session Service (`session.go`)

Manages user sessions:
- `Verify()` - Session token verification
- `Refresh()` - Session refresh
- `Revoke()` - Single session revocation
- `RevokeAll()` - All sessions revocation
- `List()` - List user sessions
- `GetCurrent()` - Get current session
- `Update()` - Update session metadata

### 4. User Service (`user.go`)

Handles user management:
- `Get()` - Get user by ID
- `GetByEmail()` - Get user by email
- `Update()` - Update user information
- `Delete()` - Delete user
- `List()` - List users with pagination
- `ChangePassword()` - Change user password
- `VerifyEmail()` - Verify user email
- `ResendVerificationEmail()` - Resend verification email
- `GetAccounts()` - Get linked OAuth accounts
- `LinkAccount()` - Link OAuth account
- `UnlinkAccount()` - Unlink OAuth account

### 5. Middleware (`middleware.go`)

HTTP middleware for web applications:
- `Authenticate()` - Require authentication
- `RequireAuth()` - Alias for Authenticate
- `OptionalAuth()` - Optional authentication
- `RequireEmailVerified()` - Require verified email
- `AuthHandler()` - Handler wrapper with auth
- `GetUserFromContext()` - Extract user from context
- `GetSessionFromContext()` - Extract session from context
- `WithTokenExtractor()` - Custom token extraction

### 6. Error Handling (`errors.go`)

Comprehensive error types:
- `Error` - Main error type
- `ErrorType` - Error type enumeration
- `ErrorResponse` - API error response
- Helper functions: `IsUnauthorizedError()`, `IsValidationError()`, etc.

### 7. Types (`types.go`)

Core data structures:
- `User` - User model
- `Session` - Session model
- `SignUpRequest/Response` - Registration types
- `SignInRequest/Response` - Login types
- `UpdateUserRequest` - User update type
- `Provider` - OAuth provider enum
- `OAuthCallbackRequest` - OAuth callback type
- `TwoFactorSetupRequest/Response` - 2FA types
- `ListUsersOptions/Response` - User listing types

### 8. Validation (`validation.go`)

Input validation utilities:
- `ValidateSignUpRequest()` - Sign up validation
- `ValidateSignInRequest()` - Sign in validation
- `ValidateEmailField()` - Email validation
- `ValidatePasswordField()` - Password strength validation
- `ValidatePasswordWithOptions()` - Custom password validation
- `ValidateNameField()` - Name validation
- `ValidateToken()` - Token validation
- `ValidateProvider()` - OAuth provider validation
- `ValidateMetadata()` - Metadata validation

### 9. Utilities (`utils.go`)

Helper functions:
- `ValidateEmail()` - Basic email validation
- `ValidatePassword()` - Password strength check
- `IsSessionExpired()` - Session expiry check
- `TimeUntilExpiry()` - Time until session expires
- `GenerateState()` - Generate OAuth state
- `SignRequest()` - HMAC request signing
- `MaskEmail()` - Mask email for display
- `MaskToken()` - Mask token for logging
- `FormatDuration()` - Human-readable duration
- `MergeMetadata()` - Merge metadata maps
- `IsValidProvider()` - Provider validation
- `GetProviderDisplayName()` - Provider display name

## Configuration

### Config Structure (`config.go`)

```go
type Config struct {
    BaseURL    string           // Better Auth server URL
    APIKey     string           // API key for authentication
    SecretKey  string           // Secret key for signing
    Timeout    time.Duration    // HTTP client timeout
    HTTPClient *http.Client     // Custom HTTP client
    Debug      bool             // Enable debug logging
}
```

## Testing

### Test Files
- `client_test.go` - Client and HTTP request tests
- `errors_test.go` - Error handling tests

### Test Coverage
- Unit tests for all core functionality
- HTTP mock server tests
- Error scenario testing
- Validation testing

### Running Tests
```bash
make test              # Run all tests
make test-coverage     # Generate coverage report
make test-verbose      # Verbose test output
```

## Examples

### 1. Basic Auth (`examples/basic_auth/`)
Demonstrates fundamental authentication flows including sign up, sign in, session verification, and user management.

### 2. Session Management (`examples/session_management/`)
Shows comprehensive session handling including verification, refresh, listing, and revocation.

### 3. HTTP Middleware (`examples/middleware/`)
Demonstrates integration with Go HTTP servers using the provided middleware.

### 4. Complete Example (`examples/complete/`)
A comprehensive example showcasing all SDK features in one application.

## Build System

### Makefile Targets
- `make build` - Build the project
- `make test` - Run tests
- `make test-coverage` - Generate coverage report
- `make lint` - Run linter
- `make fmt` - Format code
- `make vet` - Run go vet
- `make clean` - Clean artifacts
- `make install` - Install dependencies
- `make run-examples` - Run example programs
- `make all` - Run all checks

## Documentation

### Available Docs
- `README.md` - Main project documentation
- `QUICKSTART.md` - Quick start guide
- `CHANGELOG.md` - Version history
- `CONTRIBUTING.md` - Contributing guidelines
- `examples/README.md` - Examples documentation
- `doc.go` - Go package documentation

### Generating Docs
```bash
go doc -all github.com/medapsis/better-auth-sdk-go
```

## CI/CD

### GitHub Actions
- Automated testing on push/PR
- Multi-OS testing (Ubuntu, macOS, Windows)
- Multi-Go version testing (1.21, 1.22)
- Code coverage reporting
- Linting and build verification

## Dependencies

### External Dependencies
- `github.com/golang-jwt/jwt/v5` - JWT token handling

### Standard Library Usage
- `net/http` - HTTP client
- `encoding/json` - JSON encoding/decoding
- `context` - Context support
- `time` - Time operations
- `crypto/hmac` - HMAC signing
- `regexp` - Regular expressions

## Code Style

### Guidelines
- Follow Effective Go practices
- Use gofmt for formatting
- Comprehensive error handling
- Context-aware APIs
- Type-safe interfaces
- Detailed documentation
- Test coverage > 80%

### Linting
- golangci-lint configuration in `.golangci.yml`
- Enabled linters: errcheck, gosimple, govet, staticcheck, etc.

## Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

Current Version: **0.1.0** (Initial Release)

## License

MIT License - See [LICENSE](LICENSE) file for details.

## Support

- GitHub Issues: https://github.com/medapsis/better-auth-sdk-go/issues
- Documentation: https://github.com/medapsis/better-auth-sdk-go
- Better Auth Docs: https://www.better-auth.com/docs

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to this project.

---

**Note**: This project follows semantic versioning and is currently in initial development (0.x.x).