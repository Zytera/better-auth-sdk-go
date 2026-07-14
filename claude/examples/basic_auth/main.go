// Example: the two ways a backend authenticates against Better Auth — a
// session cookie or a bearer token — both resolved via the session plugin.
//
// Run: go run ./claude/examples/basic_auth
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/session"
)

var config = &betterauth.Config{BaseURL: "https://your-app.com"}

func main() {
	ctx := context.Background()

	// 1. Cookie-based: forward the end user's session cookie.
	fmt.Println("=== Cookie auth ===")
	cookieClient := betterauth.NewClient(config, &betterauth.SessionToken{
		Cookie: &http.Cookie{Name: "better-auth.session_token", Value: "session-token"},
	})
	whoAmI(ctx, cookieClient)

	// 2. Bearer-based: for token backends (mobile, service-to-service).
	fmt.Println("\n=== Bearer auth ===")
	bearerClient := betterauth.NewClient(config, &betterauth.SessionToken{})
	bearerClient.SetBearerToken("your-jwt")
	whoAmI(ctx, bearerClient)
}

func whoAmI(ctx context.Context, client *betterauth.Client) {
	data, err := session.New(client).Get(ctx)
	if err != nil {
		log.Printf("not authenticated: %v", err)
		return
	}
	fmt.Printf("Authenticated as %s (%s)\n", data.User.Name, data.User.Email)
}
