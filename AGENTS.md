# Better Auth SDK for Go — Agent Guide

This file is written for AI coding agents working on `github.com/Zytera/better-auth-sdk-go`.
It describes the project as it currently exists; prefer the code over this document when they differ.

## Project overview

This is a Go client SDK for [Better Auth](https://www.better-auth.com/). It is intentionally thin:

- The root package provides a small HTTP transport client, shared types, and typed errors.
- Every feature lives in a subpackage under `plugins/`, one per server plugin.
- Plugins depend only on the `betterauth.Requester` interface, not on the concrete `*Client`.

The SDK uses only the Go standard library (no external dependencies).

- **Module path:** `github.com/Zytera/better-auth-sdk-go`
- **Go version required:** 1.21 or later (`go.mod`)
- **License:** GNU General Public License v3.0 (see `LICENSE`)

## Project structure

```
better-auth-sdk-go/
├── client.go       # Core Client, Requester interface, bearer-token support
├── config.go       # Client configuration and defaults
├── errors.go       # Typed Better Auth errors and helpers
├── types.go        # Shared types: User, Session, SessionData
├── validation.go   # ValidationError helper
├── go.mod          # Module definition
├── README.md       # User-facing documentation
├── DEVELOPMENT.md  # How to write a new plugin
├── LICENSE         # GPL-3.0
└── plugins/        # One package per server plugin
    ├── admin/
    │   └── admin.go
    ├── session/
    │   └── session.go
    └── tenancy/
        ├── tenancy.go
        ├── tenancy_test.go
        └── types.go
```

The root package name is `betterauth`. Plugin packages have their own names (`admin`, `session`, `tenancy`) and import the root module.

## Technology stack

- Language: Go 1.21+
- HTTP client: `net/http` (custom `*http.Client` can be injected via `Config.HTTPClient`)
- JSON: `encoding/json`
- Context: all plugin methods accept `context.Context`
- Dependencies: none outside the standard library

## Build and test commands

Build the entire module:

```bash
go build ./...
```

Run all tests:

```bash
go test ./...
```

Run tests with race detection and coverage:

```bash
go test -race -cover ./...
```

Format and vet:

```bash
go fmt ./...
go vet ./...
```

There is no `Makefile` or CI workflow in the repository root at the moment; use the standard Go toolchain directly.

## Core architecture

### Client

`client.go` defines `Client` and the `Requester` interface:

```go
type Requester interface {
    Do(ctx context.Context, method, path string, body, result interface{}) error
}
```

`*Client` implements `Requester`. `Do`:

1. Marshals `body` to JSON (if non-nil).
2. Builds the URL as `Config.BaseURL + Config.BasePath + path`.
3. Adds `Content-Type: application/json` and `Accept: application/json` headers.
4. Adds the session cookie from `SessionToken.Cookie`.
5. Adds `Authorization: Bearer <token>` if `SetBearerToken` was called.
6. Sends the request and unmarshals the JSON response into `result`.
7. Returns a typed `*betterauth.Error` for HTTP >= 400.

`NewClient(config, sessionToken)` creates a client. `Config.setDefaults()` applies:

- `BasePath`: `"/api/auth"`
- `Timeout`: `30s`
- `HTTPClient`: `&http.Client{Timeout: c.Timeout}`

Plugin paths are relative to `Config.BasePath`.

### Configuration

```go
type Config struct {
    BaseURL    string
    BasePath   string
    Timeout    time.Duration
    HTTPClient *http.Client
    Debug      bool
}
```

Only `BaseURL` is semantically required. Call `config.Validate()` to check it explicitly.

### Error handling

`errors.go` defines:

- `ErrorType` constants: `validation`, `unauthorized`, `not_found`, `forbidden`, `conflict`, `internal`, `network`, `timeout`.
- `Error` struct with `Type`, `Message`, `StatusCode`, `Details`, and `Err`.
- Constructors: `NewError`, `NewErrorWithDetails`, `WrapError`.
- Type-check helpers: `IsValidationError`, `IsUnauthorizedError`, etc.

Non-2xx responses from the server are parsed by `parseErrorResponse` and turned into `*Error` values.

## Plugin conventions

Read `DEVELOPMENT.md` for the full guide. The short version:

- One package per server plugin, named after the plugin (`admin`, `session`, `tenancy`, …).
- Constructor signature: `New(r betterauth.Requester) *Plugin`.
- Accept the `Requester` interface, never `*betterauth.Client`.
- Routes are relative to `Config.BasePath`.
- Reuse core types (`betterauth.User`, `betterauth.Session`, `betterauth.SessionData`) when the server returns them.
- Validate inputs client-side and return `betterauth.NewError(betterauth.ErrorTypeValidation, …)`.
- Return server errors unchanged; `Do` already produces typed errors.
- If a route or payload is guessed rather than verified against a real server, mark it with a `// ponytail: …` comment.

### Existing plugins

| Plugin   | Import path                                  | Notes |
|----------|----------------------------------------------|-------|
| `admin`  | `github.com/Zytera/better-auth-sdk-go/plugins/admin` | User/role/impersonation admin endpoints. |
| `session`| `github.com/Zytera/better-auth-sdk-go/plugins/session` | `Get` and `Verify`. |
| `tenancy`| `github.com/Zytera/better-auth-sdk-go/plugins/tenancy` | Organizations, teams, roles, statements, members, permissions, invitations. Mirrors the TypeScript client structure. |

The bearer plugin is not a separate package; enable it with `client.SetBearerToken("jwt")`.

## Testing instructions

- Test plugins with `httptest.Server` and a real `*betterauth.Client` pointed at it.
- See `plugins/tenancy/tenancy_test.go` for the canonical pattern.
- Keep tests in `_test.go` files with package names ending in `_test` when they only exercise the public API.
- Run `go test ./...` before committing.

## Code style guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go.html).
- Run `go fmt` on all changed files.
- Run `go vet ./...` before submitting changes.
- Document all exported types and functions with godoc comments.
- Keep the root package small; add new features as plugins under `plugins/`.
- Prefer explicit, typed structs for plugin inputs/outputs over `map[string]interface{}`.
- Do not introduce external dependencies without a strong reason; the project currently has zero.

## Security considerations

- The SDK sends session cookies and bearer tokens over HTTP. In production the server must use HTTPS.
- `SessionToken.Cookie` is copied by value into requests; the caller is responsible for keeping the cookie value secret.
- `SetBearerToken` stores the token in memory on the client; do not log the client or share it between untrusted contexts.
- The SDK does not validate TLS certificates itself; it relies on `net/http` defaults and any custom `HTTPClient` provided.
- Input validation in the SDK is only for early client-side failures; the server is the authority for authentication and authorization.

## Common gotchas

- `Requester.Do` prepends `Config.BasePath` to the path you pass. Plugins must call `/session/verify`, not `/api/auth/session/verify`.
- `Config.Debug` exists but is not wired to any logger yet; adding debug logging is a future concern.
- The root package intentionally does not contain high-level services like `client.Auth` or `client.User`; use the plugin packages instead.
- The project currently has no root-level `Makefile`, `.golangci.yml`, or GitHub Actions workflow.
