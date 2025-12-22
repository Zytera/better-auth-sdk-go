package main

import (
	"context"
	"fmt"
	"log"
	"time"

	betterauth "github.com/medapsis/better-auth-sdk-go"
)

func main() {
	// Initialize the Better Auth client
	client := betterauth.NewClient(&betterauth.Config{
		BaseURL:   "https://your-app.com",
		APIKey:    "your-api-key",
		SecretKey: "your-secret-key",
		Timeout:   30 * time.Second,
	})

	ctx := context.Background()

	// Example 1: Sign in and get a session
	fmt.Println("=== Example 1: Sign In ===")
	signInResp, err := client.Auth.SignIn(ctx, &betterauth.SignInRequest{
		Email:    "user@example.com",
		Password: "securePassword123",
	})
	if err != nil {
		log.Fatalf("Failed to sign in: %v", err)
	}

	fmt.Printf("Successfully signed in!\n")
	fmt.Printf("Session ID: %s\n", signInResp.Session.ID)
	fmt.Printf("Token: %s\n", signInResp.Session.Token)
	fmt.Printf("Expires at: %s\n\n", signInResp.Session.ExpiresAt)

	sessionToken := signInResp.Session.Token
	refreshToken := signInResp.Session.RefreshToken
	userID := signInResp.User.ID

	// Example 2: Verify the session
	fmt.Println("=== Example 2: Verify Session ===")
	session, err := client.Session.Verify(ctx, sessionToken)
	if err != nil {
		log.Fatalf("Failed to verify session: %v", err)
	}

	fmt.Printf("Session is valid!\n")
	fmt.Printf("User ID: %s\n", session.UserID)
	fmt.Printf("Session expires: %s\n\n", session.ExpiresAt)

	// Example 3: Get current session details
	fmt.Println("=== Example 3: Get Current Session ===")
	currentSession, err := client.Session.GetCurrent(ctx, sessionToken)
	if err != nil {
		log.Fatalf("Failed to get current session: %v", err)
	}

	fmt.Printf("Current session ID: %s\n", currentSession.ID)
	fmt.Printf("IP Address: %s\n", currentSession.IPAddress)
	fmt.Printf("User Agent: %s\n\n", currentSession.UserAgent)

	// Example 4: List all sessions for the user
	fmt.Println("=== Example 4: List All Sessions ===")
	sessionsResp, err := client.Session.List(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to list sessions: %v", err)
	}

	fmt.Printf("Total sessions: %d\n", sessionsResp.Total)
	for i, s := range sessionsResp.Sessions {
		fmt.Printf("Session %d:\n", i+1)
		fmt.Printf("  ID: %s\n", s.ID)
		fmt.Printf("  Created: %s\n", s.CreatedAt)
		fmt.Printf("  Expires: %s\n", s.ExpiresAt)
		fmt.Printf("  IP: %s\n", s.IPAddress)
	}
	fmt.Println()

	// Example 5: Refresh the session (when it's about to expire)
	fmt.Println("=== Example 5: Refresh Session ===")
	newSession, err := client.Session.Refresh(ctx, refreshToken)
	if err != nil {
		log.Fatalf("Failed to refresh session: %v", err)
	}

	fmt.Printf("Session refreshed successfully!\n")
	fmt.Printf("New token: %s\n", newSession.Token)
	fmt.Printf("New expiry: %s\n\n", newSession.ExpiresAt)

	// Update token for subsequent requests
	sessionToken = newSession.Token

	// Example 6: Update session metadata
	fmt.Println("=== Example 6: Update Session Metadata ===")
	updates := map[string]interface{}{
		"lastActivity": time.Now().Unix(),
		"deviceType":   "desktop",
	}

	updatedSession, err := client.Session.Update(ctx, newSession.ID, updates)
	if err != nil {
		log.Fatalf("Failed to update session: %v", err)
	}

	fmt.Printf("Session updated successfully!\n")
	fmt.Printf("Session ID: %s\n\n", updatedSession.ID)

	// Example 7: Revoke a specific session
	fmt.Println("=== Example 7: Revoke Session ===")
	err = client.Session.Revoke(ctx, sessionToken)
	if err != nil {
		log.Fatalf("Failed to revoke session: %v", err)
	}

	fmt.Printf("Session revoked successfully!\n\n")

	// Example 8: Try to verify the revoked session (should fail)
	fmt.Println("=== Example 8: Verify Revoked Session ===")
	_, err = client.Session.Verify(ctx, sessionToken)
	if err != nil {
		if betterauth.IsUnauthorizedError(err) {
			fmt.Printf("Session is invalid (as expected): %v\n\n", err)
		} else {
			log.Fatalf("Unexpected error: %v", err)
		}
	}

	// Example 9: Sign in again for demonstration
	fmt.Println("=== Example 9: Sign In Again ===")
	signInResp2, err := client.Auth.SignIn(ctx, &betterauth.SignInRequest{
		Email:    "user@example.com",
		Password: "securePassword123",
	})
	if err != nil {
		log.Fatalf("Failed to sign in: %v", err)
	}

	fmt.Printf("Signed in again!\n")
	fmt.Printf("New session token: %s\n\n", signInResp2.Session.Token)

	// Example 10: Revoke all sessions for the user
	fmt.Println("=== Example 10: Revoke All Sessions ===")
	err = client.Session.RevokeAll(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to revoke all sessions: %v", err)
	}

	fmt.Printf("All sessions revoked successfully!\n")
	fmt.Printf("User has been signed out from all devices.\n")

	// Example 11: Error handling
	fmt.Println("\n=== Example 11: Error Handling ===")
	_, err = client.Session.Verify(ctx, "invalid-token")
	if err != nil {
		if betterauth.IsUnauthorizedError(err) {
			fmt.Println("✓ Correctly detected unauthorized error")
		} else if betterauth.IsValidationError(err) {
			fmt.Println("✓ Correctly detected validation error")
		} else {
			fmt.Printf("Other error: %v\n", err)
		}
	}

	fmt.Println("\n=== Session Management Examples Complete ===")
}
