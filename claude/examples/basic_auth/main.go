package main

import (
	"context"
	"fmt"
	"log"
	"time"

	betterauth "github.com/Zytera/better-auth-sdk-go"
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

	// Example 1: Sign up a new user
	fmt.Println("=== Sign Up Example ===")
	signUpResp, err := client.Auth.SignUp(ctx, &betterauth.SignUpRequest{
		Email:    "john.doe@example.com",
		Password: "SecurePassword123!",
		Name:     "John Doe",
	})
	if err != nil {
		if betterauth.IsConflictError(err) {
			log.Printf("User already exists: %v\n", err)
		} else {
			log.Fatalf("Failed to sign up: %v", err)
		}
	} else {
		fmt.Printf("✓ User created successfully!\n")
		fmt.Printf("  User ID: %s\n", signUpResp.User.ID)
		fmt.Printf("  Email: %s\n", signUpResp.User.Email)
		fmt.Printf("  Name: %s\n", signUpResp.User.Name)
		fmt.Printf("  Session Token: %s\n", signUpResp.Session.Token)
		fmt.Println()
	}

	// Example 2: Sign in with email and password
	fmt.Println("=== Sign In Example ===")
	signInResp, err := client.Auth.SignIn(ctx, &betterauth.SignInRequest{
		Email:    "john.doe@example.com",
		Password: "SecurePassword123!",
	})
	if err != nil {
		if betterauth.IsUnauthorizedError(err) {
			log.Printf("Invalid credentials: %v\n", err)
		} else {
			log.Fatalf("Failed to sign in: %v", err)
		}
	} else {
		fmt.Printf("✓ User signed in successfully!\n")
		fmt.Printf("  User ID: %s\n", signInResp.User.ID)
		fmt.Printf("  Email: %s\n", signInResp.User.Email)
		fmt.Printf("  Session Token: %s\n", signInResp.Session.Token)
		fmt.Printf("  Session expires at: %s\n", signInResp.Session.ExpiresAt)
		fmt.Println()
	}

	// Store session token for subsequent requests
	sessionToken := signInResp.Session.Token

	// Example 3: Verify the session
	fmt.Println("=== Verify Session Example ===")
	session, err := client.Session.Verify(ctx, sessionToken)
	if err != nil {
		log.Fatalf("Failed to verify session: %v", err)
	}
	fmt.Printf("✓ Session is valid!\n")
	fmt.Printf("  Session ID: %s\n", session.ID)
	fmt.Printf("  User ID: %s\n", session.UserID)
	fmt.Printf("  Expires at: %s\n", session.ExpiresAt)
	fmt.Println()

	// Example 4: Get user information
	fmt.Println("=== Get User Example ===")
	user, err := client.User.Get(ctx, signInResp.User.ID)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	fmt.Printf("✓ User retrieved successfully!\n")
	fmt.Printf("  ID: %s\n", user.ID)
	fmt.Printf("  Email: %s\n", user.Email)
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Email Verified: %t\n", user.EmailVerified)
	fmt.Printf("  Created At: %s\n", user.CreatedAt)
	fmt.Println()

	// Example 5: Update user information
	fmt.Println("=== Update User Example ===")
	updatedUser, err := client.User.Update(ctx, user.ID, &betterauth.UpdateUserRequest{
		Name: "John Updated Doe",
	})
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	fmt.Printf("✓ User updated successfully!\n")
	fmt.Printf("  New Name: %s\n", updatedUser.Name)
	fmt.Println()

	// Example 6: Change password
	fmt.Println("=== Change Password Example ===")
	err = client.Auth.ChangePassword(ctx, sessionToken, &betterauth.ChangePasswordRequest{
		CurrentPassword: "SecurePassword123!",
		NewPassword:     "NewSecurePassword456!",
	})
	if err != nil {
		if betterauth.IsValidationError(err) {
			log.Printf("Validation error: %v\n", err)
		} else {
			log.Printf("Failed to change password: %v", err)
		}
	} else {
		fmt.Printf("✓ Password changed successfully!\n")
		fmt.Println()
	}

	// Example 7: Sign out
	fmt.Println("=== Sign Out Example ===")
	err = client.Auth.SignOut(ctx, sessionToken)
	if err != nil {
		log.Fatalf("Failed to sign out: %v", err)
	}
	fmt.Printf("✓ User signed out successfully!\n")
	fmt.Println()

	// Example 8: Verify session after sign out (should fail)
	fmt.Println("=== Verify Session After Sign Out ===")
	_, err = client.Session.Verify(ctx, sessionToken)
	if err != nil {
		if betterauth.IsUnauthorizedError(err) {
			fmt.Printf("✓ Session is invalid as expected: %v\n", err)
		} else {
			log.Printf("Unexpected error: %v", err)
		}
	} else {
		log.Printf("Warning: Session should be invalid after sign out")
	}

	fmt.Println("\n=== All examples completed! ===")
}
