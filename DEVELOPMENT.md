# Writing a plugin

The SDK mirrors [Better Auth](https://www.better-auth.com/)'s plugin model: the
core `betterauth.Client` is a thin HTTP transport, and every feature lives in
its own subpackage under `plugins/`. A plugin never touches the core — it just
accepts the client through one small interface and calls the server.

## The contract

The only thing a plugin needs from the client is `Requester`:

```go
// package betterauth
type Requester interface {
    Do(ctx context.Context, method, path string, body, result interface{}) error
}
```

`*betterauth.Client` implements it. `Do` prepends the configured base path
(`Config.BasePath`, default `/api/auth`) to `path`, marshals `body` to JSON,
sends the session cookie (and `Authorization: Bearer` if set), and unmarshals
the response into `result` (pass `nil` when you don't need the body). Non-2xx
responses come back as a typed `*betterauth.Error`.

**Plugin paths are relative to the base path** — pass `/my-plugin/do-thing`,
not `/api/auth/my-plugin/do-thing`. The client adds the prefix, so a custom
server `basePath` just works.

Because plugins depend only on this interface — not on the concrete client —
they are trivially testable with a mock and can live in **any** module,
including a third-party repo.

## Anatomy of a plugin

Every plugin follows the same shape. Here's the real `admin` plugin
(`plugins/admin/admin.go`), trimmed to two endpoints — copy this and change the
routes/payloads:

```go
// Package admin is the client-side plugin for the better-auth admin plugin.
package admin

import (
    "context"

    betterauth "github.com/Zytera/better-auth-sdk-go"
)

type Plugin struct {
    r betterauth.Requester
}

// New wires the plugin to any Requester (typically *betterauth.Client).
func New(r betterauth.Requester) *Plugin {
    return &Plugin{r: r}
}

// SetRole sets a user's role. Returns only an error (no response body needed).
func (p *Plugin) SetRole(ctx context.Context, userID, role string) error {
    return p.r.Do(ctx, "POST", "/admin/set-role", map[string]string{
        "userId": userID,
        "role":   role,
    }, nil)
}

// CreateUser creates a user and decodes the returned user object.
func (p *Plugin) CreateUser(ctx context.Context, in CreateUserInput) (*betterauth.User, error) {
    var out struct {
        User betterauth.User `json:"user"`
    }
    if err := p.r.Do(ctx, "POST", "/admin/create-user", in, &out); err != nil {
        return nil, err
    }
    return &out.User, nil
}
```

## Using the admin plugin

Build the client once, then construct the plugin with it:

```go
package main

import (
    "context"
    "log"
    "net/http"

    betterauth "github.com/Zytera/better-auth-sdk-go"
    "github.com/Zytera/better-auth-sdk-go/plugins/admin"
)

func main() {
    // 1. Core client, authenticated as an admin user (cookie or bearer token).
    client := betterauth.NewClient(
        &betterauth.Config{BaseURL: "https://your-app.com"},
        &betterauth.SessionToken{Cookie: &http.Cookie{
            Name:  "better-auth.session_token",
            Value: "admin-session-token",
        }},
    )
    // (or, for a token-based backend: client.SetBearerToken("admin-jwt"))

    // 2. The admin plugin — just pass it the client.
    adminClient := admin.New(client)

    ctx := context.Background()

    // 3. Call endpoints.
    user, err := adminClient.CreateUser(ctx, admin.CreateUserInput{
        Email:    "user@example.com",
        Password: "s3cret!",
        Name:     "New User",
        Role:     "user",
    })
    if err != nil {
        log.Fatal(err)
    }

    if err := adminClient.SetRole(ctx, user.ID, "admin"); err != nil {
        log.Fatal(err)
    }
    log.Printf("promoted %s to admin", user.Email)
}
```

Every other plugin is used the exact same way: `pkg.New(client)`, then call.

## Conventions

- **One package per server plugin**, named after it (lowercase, no separators:
  `phonenumber`, `expopasskey`).
- **`New(r betterauth.Requester) *Plugin`** is the constructor. Always accept
  the interface, never the concrete `*Client`.
- **Paths are relative** to `Config.BasePath` (default `/api/auth`). A server
  plugin registered at `myPlugin` exposes routes under `/my-plugin/...`.
- **Reuse core types** (`betterauth.User`, `betterauth.Session`,
  `betterauth.SessionData`) when the endpoint returns them; define
  plugin-specific structs otherwise.
- **Validate required inputs** before the request and return
  `betterauth.NewError(betterauth.ErrorTypeValidation, ...)`.
- **Return errors as-is** — `Do` already produces typed `*betterauth.Error`s.
- **Flag unverified routes.** If you can't confirm a path/payload against the
  server yet, leave a `// ponytail: route guessed — confirm against server`
  comment so it's not mistaken for verified.

## Third-party plugins

A plugin doesn't have to live in this repo. In your own module:

```go
import betterauth "github.com/Zytera/better-auth-sdk-go"

type Plugin struct{ r betterauth.Requester }
func New(r betterauth.Requester) *Plugin { return &Plugin{r} }
```

Users pass the same `*betterauth.Client` to your `New`. No registration, no
core changes.

## Testing

Spin up an `httptest.Server`, point the client at it, and assert the request.
See `plugins/phonenumber/phonenumber_test.go` for the pattern:

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // assert r.URL.Path / decode r.Body, then write a JSON response
}))
defer srv.Close()

c := betterauth.NewClient(&betterauth.Config{BaseURL: srv.URL}, &betterauth.SessionToken{})
p := admin.New(c)
// ... call p and assert
```

Run everything with:

```bash
go build ./... && go test ./plugins/...
```
