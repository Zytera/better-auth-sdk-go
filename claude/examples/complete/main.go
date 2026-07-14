// Example: a tour of the plugin API — session, admin, tenancy.
//
// Each plugin is constructed the same way: pkg.New(client).
// Run: go run ./claude/examples/complete
package main

import (
	"context"
	"fmt"
	"log"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/admin"
	"github.com/Zytera/better-auth-sdk-go/plugins/session"
	"github.com/Zytera/better-auth-sdk-go/plugins/tenancy"
)

func main() {
	// Authenticated as an admin (bearer token) for the whole tour.
	client := betterauth.NewClient(
		&betterauth.Config{BaseURL: "https://your-app.com"},
		&betterauth.SessionToken{},
	)
	client.SetBearerToken("admin-jwt-token")
	ctx := context.Background()

	// --- read the current session -------------------------------------------
	sess := session.New(client)
	if cur, err := sess.Get(ctx); err == nil {
		fmt.Printf("Active session %s for %s\n", cur.Session.ID, cur.User.Email)
	}

	// --- admin operations ---------------------------------------------------
	adm := admin.New(client)
	user, err := adm.CreateUser(ctx, admin.CreateUserInput{
		Email: "user@example.com", Password: "s3cret!", Name: "New User", Role: "user",
	})
	if err != nil {
		log.Printf("create user: %v", err)
	} else {
		_ = adm.SetRole(ctx, user.ID, "admin")
		fmt.Printf("Created + promoted %s\n", user.Email)
	}

	// --- tenancy: org -> team -> member -------------------------------------
	ten := tenancy.New(client)
	org, err := ten.Organization.Create(ctx, tenancy.CreateOrgInput{Name: "Acme", Slug: "acme"})
	if err != nil {
		log.Printf("create org: %v", err)
		return
	}
	team, _ := ten.Team.Create(ctx, tenancy.CreateTeamInput{
		ParentType: tenancy.ContextOrganization, ParentID: org.ID,
		Name: "Engineering", Slug: "engineering",
	})
	if user != nil {
		ten.Member.Add(ctx, tenancy.AddMemberInput{
			UserID: user.ID, ContextType: tenancy.ContextTeam, ContextID: team.ID,
		})
	}
	fmt.Printf("Org %s / Team %s ready\n", org.ID, team.ID)
}
