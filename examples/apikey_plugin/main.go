// Example: creating an API key, then authenticating calls with it.
//
// Run: go run ./examples/apikey_plugin
package main

import (
	"context"
	"fmt"
	"log"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/apikey"
	"github.com/Zytera/better-auth-sdk-go/plugins/session"
)

func main() {
	// 1. Core client, authenticated as a user (bearer here; a cookie works too).
	client := betterauth.NewClient(
		&betterauth.Config{BaseURL: "https://your-app.com"},
		&betterauth.SessionToken{},
	)
	client.SetBearerToken("user-jwt-token")

	ctx := context.Background()

	// 2. Create an API key. The plaintext key is returned exactly once.
	key, err := apikey.New(client).Create(ctx, apikey.CreateInput{
		Name:   "ci-pipeline",
		Prefix: "ba_",
	})
	if err != nil {
		log.Fatalf("create key: %v", err)
	}
	fmt.Printf("Created key %s: %s\n", key.ID, key.Key)

	// 3. Authenticate subsequent calls *with* that key — sent as x-api-key.
	//    Any endpoint now resolves to the key's owner, no session needed.
	client.SetAPIKey(key.Key)

	data, err := session.New(client).Get(ctx)
	if err != nil {
		log.Fatalf("get session: %v", err)
	}
	fmt.Printf("Authenticated as %s via API key\n", data.User.Email)
}
