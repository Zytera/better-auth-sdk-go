// Example: session management via the session plugin.
//
// Run: go run ./examples/session_management
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
	// A client is bound to one session token (here, a cookie). To act as a
	// different user, build another client with that user's token.
	client := betterauth.NewClient(
		&betterauth.Config{BaseURL: "https://your-app.com"},
		&betterauth.SessionToken{Cookie: &http.Cookie{
			Name:  "better-auth.session_token",
			Value: "your-session-token-here",
		}},
	)
	sess := session.New(client)

	ctx := context.Background()

	// Current session + user for this client's token.
	fmt.Println("=== Get Session ===")
	data, err := sess.Get(ctx)
	if err != nil {
		log.Fatalf("get session: %v", err)
	}
	fmt.Printf("User:    %s (%s)\n", data.User.Name, data.User.Email)
	fmt.Printf("Session: %s (expires %s)\n\n", data.Session.ID, data.Session.ExpiresAt)

	// Verify an arbitrary token (e.g. one received from a client app).
	fmt.Println("=== Verify Token ===")
	s, err := sess.Verify(ctx, data.Session.Token)
	if err != nil {
		log.Fatalf("verify: %v", err)
	}
	fmt.Printf("Token valid for user %s until %s\n", s.UserID, s.ExpiresAt)

	// Note: listing/revoking sessions is not wrapped in the session plugin yet.
}
