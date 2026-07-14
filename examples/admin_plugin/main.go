// Example: using the admin plugin.
//
// The client must be authenticated as an admin user (cookie or bearer token).
// Run: go run ./examples/admin_plugin
package main

import (
	"context"
	"fmt"
	"log"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/admin"
)

func main() {
	// 1. Core client. Here we authenticate with a bearer token (e.g. a mobile
	//    backend); a session cookie works too via SessionToken.
	client := betterauth.NewClient(
		&betterauth.Config{BaseURL: "https://your-app.com"},
		&betterauth.SessionToken{},
	)
	client.SetBearerToken("admin-jwt-token")

	// 2. The admin plugin.
	adm := admin.New(client)

	ctx := context.Background()

	// 3. Create a user with an explicit role.
	user, err := adm.CreateUser(ctx, admin.CreateUserInput{
		Email:    "user@example.com",
		Password: "s3cret!",
		Name:     "New User",
		Role:     "user",
	})
	if err != nil {
		log.Fatalf("create user: %v", err)
	}
	fmt.Printf("Created user %s (%s)\n", user.ID, user.Email)

	// 4. Promote to admin.
	if err := adm.SetRole(ctx, user.ID, "admin"); err != nil {
		log.Fatalf("set role: %v", err)
	}
	fmt.Println("Promoted to admin")

	// 5. Ban / unban.
	if err := adm.BanUser(ctx, user.ID, "spam", 3600); err != nil {
		log.Fatalf("ban: %v", err)
	}
	fmt.Println("Banned for 1h")
	if err := adm.UnbanUser(ctx, user.ID); err != nil {
		log.Fatalf("unban: %v", err)
	}

	// 6. List users.
	list, err := adm.ListUsers(ctx, admin.ListUsersQuery{Limit: 10})
	if err != nil {
		log.Fatalf("list: %v", err)
	}
	fmt.Printf("Total users: %d\n", list.Total)

	// 7. Impersonate then stop.
	if _, err := adm.ImpersonateUser(ctx, user.ID); err == nil {
		_ = adm.StopImpersonating(ctx)
	}
}
