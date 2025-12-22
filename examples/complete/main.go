package main

import (
	"context"
	"fmt"
	"log"
	"time"

	betterauth "github.com/medapsis/better-auth-sdk-go"
)

func main() {
	fmt.Println("=== Better Auth SDK - Complete Example ===")
	fmt.Println()

	// Initialize the Better Auth client
	client := betterauth.NewClient(&betterauth.Config{
		BaseURL:   "https://your-app.com",
		APIKey:    "your-api-key",
		SecretKey: "your-secret-key",
		Timeout:   30 * time.Second,
	})

	ctx := context.Background()

	// Example 1: User Registration
	fmt.Println("1. User Registration")
	fmt.Println("-------------------")
	signUpResp, err := client.Auth.SignUp(ctx, &betterauth.SignUpRequest{
		Email:    "demo@example.com",
		Password: "SecurePassword123!",
		Name:     "Demo User",
		Metadata: map[string]interface{}{
			"role":       "user",
			"department": "engineering",
		},
	})
	if err != nil {
		if betterauth.IsConflictError(err) {
			fmt.Println("✓ User already exists (expected for demo)")
		} else {
			log.Printf("Error signing up: %v\n", err)
		}
	} else {
		fmt.Printf("✓ User created successfully\n")
		fmt.Printf("  User ID: %s\n", signUpResp.User.ID)
		fmt.Printf("  Email: %s\n", signUpResp.User.Email)
		fmt.Printf("  Name: %s\n", signUpResp.User.Name)
	}
	fmt.Println()

	// Example 2: User Authentication
	fmt.Println("2. User Authentication")
	fmt.Println("---------------------")
	signInResp, err := client.Auth.SignIn(ctx, &betterauth.SignInRequest{
		Email:    "demo@example.com",
		Password: "SecurePassword123!",
	})
	if err != nil {
		log.Fatalf("Failed to sign in: %v", err)
	}
	fmt.Printf("✓ User authenticated successfully\n")
	fmt.Printf("  Session Token: %s...\n", betterauth.MaskToken(signInResp.Session.Token))
	fmt.Printf("  Expires: %s\n", signInResp.Session.ExpiresAt.Format(time.RFC3339))
	fmt.Println()

	sessionToken := signInResp.Session.Token
	userID := signInResp.User.ID

	// Example 3: Session Verification
	fmt.Println("3. Session Verification")
	fmt.Println("----------------------")
	session, err := client.Session.Verify(ctx, sessionToken)
	if err != nil {
		log.Fatalf("Failed to verify session: %v", err)
	}
	fmt.Printf("✓ Session verified successfully\n")
	fmt.Printf("  Session ID: %s\n", session.ID)
	fmt.Printf("  User ID: %s\n", session.UserID)
	fmt.Printf("  Time until expiry: %s\n", betterauth.FormatDuration(betterauth.TimeUntilExpiry(session.ExpiresAt)))
	fmt.Println()

	// Example 4: Get User Profile
	fmt.Println("4. Get User Profile")
	fmt.Println("------------------")
	user, err := client.User.Get(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	fmt.Printf("✓ User profile retrieved\n")
	fmt.Printf("  ID: %s\n", user.ID)
	fmt.Printf("  Email: %s\n", user.Email)
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Email Verified: %t\n", user.EmailVerified)
	fmt.Printf("  Created: %s\n", user.CreatedAt.Format("2006-01-02"))
	fmt.Println()

	// Example 5: Update User Profile
	fmt.Println("5. Update User Profile")
	fmt.Println("---------------------")
	updatedUser, err := client.User.Update(ctx, userID, &betterauth.UpdateUserRequest{
		Name: "Demo User (Updated)",
		Metadata: map[string]interface{}{
			"role":       "admin",
			"department": "engineering",
			"updated":    time.Now().Unix(),
		},
	})
	if err != nil {
		log.Printf("Failed to update user: %v\n", err)
	} else {
		fmt.Printf("✓ User profile updated\n")
		fmt.Printf("  New Name: %s\n", updatedUser.Name)
	}
	fmt.Println()

	// Example 6: List All Sessions
	fmt.Println("6. List User Sessions")
	fmt.Println("--------------------")
	sessionsResp, err := client.Session.List(ctx, userID)
	if err != nil {
		log.Printf("Failed to list sessions: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d active session(s)\n", sessionsResp.Total)
		for i, s := range sessionsResp.Sessions {
			fmt.Printf("  Session %d:\n", i+1)
			fmt.Printf("    ID: %s\n", s.ID)
			fmt.Printf("    Created: %s\n", s.CreatedAt.Format(time.RFC3339))
			fmt.Printf("    IP: %s\n", s.IPAddress)
		}
	}
	fmt.Println()

	// Example 7: Password Reset Flow
	fmt.Println("7. Password Reset Flow")
	fmt.Println("---------------------")
	err = client.Auth.ResetPassword(ctx, &betterauth.ResetPasswordRequest{
		Email: "demo@example.com",
	})
	if err != nil {
		log.Printf("Failed to initiate password reset: %v\n", err)
	} else {
		fmt.Println("✓ Password reset email sent")
		fmt.Println("  Check email for reset link")
	}
	fmt.Println()

	// Example 8: OAuth URL Generation
	fmt.Println("8. OAuth Authentication")
	fmt.Println("----------------------")
	oauthState := betterauth.GenerateState()
	providers := []betterauth.Provider{
		betterauth.ProviderGoogle,
		betterauth.ProviderGithub,
		betterauth.ProviderFacebook,
	}

	for _, provider := range providers {
		urlResp, err := client.Auth.GetOAuthURL(ctx, provider, oauthState)
		if err != nil {
			log.Printf("Failed to get OAuth URL for %s: %v\n", provider, err)
		} else {
			fmt.Printf("✓ %s OAuth URL generated\n", betterauth.GetProviderDisplayName(provider))
			fmt.Printf("  URL: %s\n", urlResp.URL[:50]+"...")
		}
	}
	fmt.Println()

	// Example 9: Email Verification
	fmt.Println("9. Email Verification")
	fmt.Println("--------------------")
	err = client.Auth.SendVerificationEmail(ctx, "demo@example.com")
	if err != nil {
		log.Printf("Failed to send verification email: %v\n", err)
	} else {
		fmt.Println("✓ Verification email sent")
		fmt.Println("  Check email for verification link")
	}
	fmt.Println()

	// Example 10: Get User by Email
	fmt.Println("10. Get User by Email")
	fmt.Println("--------------------")
	userByEmail, err := client.User.GetByEmail(ctx, "demo@example.com")
	if err != nil {
		log.Printf("Failed to get user by email: %v\n", err)
	} else {
		fmt.Printf("✓ User found\n")
		fmt.Printf("  ID: %s\n", userByEmail.ID)
		fmt.Printf("  Email: %s (masked: %s)\n", userByEmail.Email, betterauth.MaskEmail(userByEmail.Email))
	}
	fmt.Println()

	// Example 11: List Users
	fmt.Println("11. List Users")
	fmt.Println("-------------")
	usersResp, err := client.User.List(ctx, &betterauth.ListUsersOptions{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		log.Printf("Failed to list users: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved %d user(s) (Total: %d)\n", len(usersResp.Users), usersResp.Total)
		for i, u := range usersResp.Users {
			if i < 3 { // Show only first 3
				fmt.Printf("  User %d: %s <%s>\n", i+1, u.Name, u.Email)
			}
		}
		if len(usersResp.Users) > 3 {
			fmt.Printf("  ... and %d more\n", len(usersResp.Users)-3)
		}
	}
	fmt.Println()

	// Example 12: Validation Utilities
	fmt.Println("12. Validation Utilities")
	fmt.Println("-----------------------")
	testEmails := []string{
		"valid@example.com",
		"invalid-email",
		"",
	}

	for _, email := range testEmails {
		if betterauth.ValidateEmail(email) {
			fmt.Printf("✓ '%s' is valid\n", email)
		} else {
			fmt.Printf("✗ '%s' is invalid\n", email)
		}
	}
	fmt.Println()

	// Example 13: Password Strength Validation
	fmt.Println("13. Password Strength")
	fmt.Println("--------------------")
	testPasswords := []string{
		"SecurePassword123!",
		"weak",
		"NoNumbers!",
	}

	for _, pwd := range testPasswords {
		err := betterauth.ValidatePassword(pwd)
		if err != nil {
			fmt.Printf("✗ '%s' - %v\n", pwd, err)
		} else {
			fmt.Printf("✓ '%s' is strong\n", pwd)
		}
	}
	fmt.Println()

	// Example 14: Session Refresh
	fmt.Println("14. Session Refresh")
	fmt.Println("------------------")
	if signInResp.Session.RefreshToken != "" {
		newSession, err := client.Session.Refresh(ctx, signInResp.Session.RefreshToken)
		if err != nil {
			log.Printf("Failed to refresh session: %v\n", err)
		} else {
			fmt.Println("✓ Session refreshed successfully")
			fmt.Printf("  New Token: %s...\n", betterauth.MaskToken(newSession.Token))
			fmt.Printf("  New Expiry: %s\n", newSession.ExpiresAt.Format(time.RFC3339))
			sessionToken = newSession.Token
		}
	} else {
		fmt.Println("ℹ Refresh token not available")
	}
	fmt.Println()

	// Example 15: Error Handling Demo
	fmt.Println("15. Error Handling")
	fmt.Println("-----------------")
	_, err = client.Session.Verify(ctx, "invalid-token")
	if err != nil {
		fmt.Printf("✓ Error detected: %v\n", err)
		if betterauth.IsUnauthorizedError(err) {
			fmt.Println("  Type: Unauthorized Error")
		} else if betterauth.IsValidationError(err) {
			fmt.Println("  Type: Validation Error")
		} else {
			fmt.Println("  Type: Other Error")
		}
	}
	fmt.Println()

	// Example 16: Metadata Operations
	fmt.Println("16. Metadata Operations")
	fmt.Println("----------------------")
	baseMetadata := map[string]interface{}{
		"theme":    "dark",
		"language": "en",
	}
	updates := map[string]interface{}{
		"theme":      "light",
		"last_login": time.Now().Unix(),
	}
	merged := betterauth.MergeMetadata(baseMetadata, updates)
	fmt.Printf("✓ Metadata merged\n")
	fmt.Printf("  Theme: %v\n", betterauth.SafeString(merged, "theme"))
	fmt.Printf("  Language: %v\n", betterauth.SafeString(merged, "language"))
	fmt.Printf("  Last Login: %v\n", betterauth.SafeInt(merged, "last_login"))
	fmt.Println()

	// Example 17: Sign Out
	fmt.Println("17. Sign Out")
	fmt.Println("-----------")
	err = client.Auth.SignOut(ctx, sessionToken)
	if err != nil {
		log.Printf("Failed to sign out: %v\n", err)
	} else {
		fmt.Println("✓ User signed out successfully")
	}
	fmt.Println()

	// Example 18: Verify Session After Sign Out
	fmt.Println("18. Verify After Sign Out")
	fmt.Println("------------------------")
	_, err = client.Session.Verify(ctx, sessionToken)
	if err != nil {
		if betterauth.IsUnauthorizedError(err) {
			fmt.Println("✓ Session is invalid (as expected)")
		} else {
			fmt.Printf("Unexpected error: %v\n", err)
		}
	} else {
		fmt.Println("⚠ Warning: Session should be invalid after sign out")
	}
	fmt.Println()

	// Summary
	fmt.Println("===========================================")
	fmt.Println("✓ All examples completed successfully!")
	fmt.Println("===========================================")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Check the middleware example for HTTP integration")
	fmt.Println("  • Read the documentation at github.com/medapsis/better-auth-sdk-go")
	fmt.Println("  • Join our community for support")
}
