# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Better Auth SDK for Go
- Core client with authentication, session, and user management
- Sign up, sign in, and sign out functionality
- Session management (verify, refresh, revoke)
- User CRUD operations
- OAuth/Social authentication support (Google, GitHub, Facebook, Twitter, Apple, Discord, Microsoft)
- Two-factor authentication (2FA) setup and verification
- Email verification support
- Password reset functionality
- HTTP middleware for authentication
- Comprehensive error handling with typed errors
- Utility functions for validation and security
- Context-based user and session retrieval
- Extensive test coverage
- Example applications demonstrating usage
- Complete documentation and API reference

### Features
- ✅ Email/Password authentication
- ✅ OAuth/Social login support
- ✅ Session management
- ✅ User management
- ✅ Two-factor authentication
- ✅ Email verification
- ✅ Password reset
- ✅ HTTP middleware
- ✅ Custom error types
- ✅ Context-based authentication
- ✅ Thread-safe client
- ✅ Configurable timeouts
- ✅ Comprehensive logging support

## [0.1.0] - 2024-01-XX

### Added
- Initial project structure
- Core SDK implementation
- Basic authentication flows
- Documentation and examples

---

## Release Notes

### Version 0.1.0 (Initial Release)

This is the initial release of the Better Auth SDK for Go. It provides a complete client library for integrating Better Auth authentication into your Go applications.

**Key Features:**
- Complete authentication flow support
- Easy to use and integrate
- Type-safe API
- Comprehensive error handling
- Session management
- User management
- OAuth/Social authentication
- Two-factor authentication
- HTTP middleware for easy integration

**Installation:**
```bash
go get github.com/Zytera/better-auth-sdk-go
```

**Quick Start:**
```go
client := betterauth.NewClient(&betterauth.Config{
    BaseURL:   "https://your-app.com",
    APIKey:    "your-api-key",
    SecretKey: "your-secret-key",
})

// Sign up
user, err := client.Auth.SignUp(ctx, &betterauth.SignUpRequest{
    Email:    "user@example.com",
    Password: "securePassword123",
    Name:     "John Doe",
})
```

For more information, see the [README](README.md) and [examples](examples/).

---

[Unreleased]: https://github.com/Zytera/better-auth-sdk-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Zytera/better-auth-sdk-go/releases/tag/v0.1.0
