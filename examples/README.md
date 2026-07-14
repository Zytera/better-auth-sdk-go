# Better Auth SDK — Examples

Runnable examples for the Go SDK. Each is its own `package main`; edit the
`BaseURL` and token placeholders before running. All use the current plugin API
(`plugins/*`).

| Example | Shows |
|---------|-------|
| `basic_auth/` | The two backend auth modes — session cookie vs bearer token — resolved with the session plugin |
| `session_plugin/` | Session basics: `Get` the current session/user, `Verify` a token |
| `session_management/` | Same, framed around managing a client's session |
| `middleware/` | An HTTP middleware that authenticates requests with the session plugin and puts the user in the context |
| `admin_plugin/` | User administration: create, set role, ban/unban, list, impersonate |
| `complete/` | A tour across session → admin → tenancy (org → team → member) |

## Run

```bash
go run ./examples/session_plugin
go run ./examples/admin_plugin
go run ./examples/middleware   # starts an HTTP server on :8080
```

## Pattern

Every example follows the same shape: build the client once, then construct
each plugin with it.

```go
client := betterauth.NewClient(cfg, sessionToken)
sess   := session.New(client)
adm    := admin.New(client)

data, _ := sess.Get(ctx)
adm.SetRole(ctx, userID, "admin")
```

See the top-level [README](../README.md) and [DEVELOPMENT.md](../DEVELOPMENT.md)
for the full plugin list and how to write your own.
