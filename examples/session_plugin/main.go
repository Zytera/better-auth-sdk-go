// Example: using the session plugin.
//
// Run: go run ./examples/session_plugin
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/session"
)

func main() {
	// 1. Core client authenticated with a session cookie.
	client := betterauth.NewClient(
		&betterauth.Config{BaseURL: "https://your-app.com"},
		&betterauth.SessionToken{Cookie: &http.Cookie{
			Name:  "better-auth.session_token",
			Value: "your-session-token-here",
		}},
	)

	// 2. The session plugin — just pass it the client.
	sess := session.New(client)

	ctx := context.Background()

	// 3a. Get the current session + user for the cookie above.
	data, err := sess.Get(ctx)
	if err != nil {
		log.Fatalf("get session: %v", err)
	}
	fmt.Printf("User:    %s (%s)\n", data.User.Name, data.User.Email)
	fmt.Printf("Session: %s (expires %s)\n", data.Session.ID, data.Session.ExpiresAt)

	// 3b. Verify an arbitrary token.
	s, err := sess.Verify(ctx, data.Session.Token)
	if err != nil {
		log.Fatalf("verify: %v", err)
	}
	fmt.Printf("Token valid for user %s\n", s.UserID)
}
