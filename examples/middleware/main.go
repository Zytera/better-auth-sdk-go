// Example: HTTP middleware that authenticates requests with the session plugin.
//
// Run: go run ./examples/middleware
// Then: curl -H 'Cookie: better-auth.session_token=TOKEN' localhost:8080/api/protected
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/session"
)

const cookieName = "better-auth.session_token"

type ctxKey string

const userKey ctxKey = "user"

var config = &betterauth.Config{BaseURL: "https://your-app.com"}

// authenticate verifies the request's session cookie and puts the user in the
// context. A client is built per request because each carries its own token.
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		client := betterauth.NewClient(config, &betterauth.SessionToken{Cookie: cookie})
		data, err := session.New(client).Get(r.Context())
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, &data.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userFrom(ctx context.Context) *betterauth.User {
	u, _ := ctx.Value(userKey).(*betterauth.User)
	return u
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/public", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"message": "public endpoint"})
	})
	mux.Handle("/api/protected", authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := userFrom(r.Context())
		writeJSON(w, map[string]any{"id": u.ID, "email": u.Email, "name": u.Name})
	})))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"status": "healthy", "time": time.Now().Format(time.RFC3339)})
	})

	fmt.Println("Server on :8080  (/health, /api/public, /api/protected)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
