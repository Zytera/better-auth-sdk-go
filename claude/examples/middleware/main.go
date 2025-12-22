package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// Create middleware instance
	middleware := betterauth.NewMiddleware(client)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Public endpoint - no authentication required
	mux.HandleFunc("/api/public", publicHandler)

	// Protected endpoint - requires authentication
	mux.Handle("/api/protected", middleware.Authenticate(http.HandlerFunc(protectedHandler)))

	// Optional auth endpoint - works with or without authentication
	mux.Handle("/api/optional", middleware.OptionalAuth(http.HandlerFunc(optionalAuthHandler)))

	// Requires verified email
	mux.Handle("/api/verified-only",
		middleware.Authenticate(
			middleware.RequireEmailVerified(http.HandlerFunc(verifiedOnlyHandler)),
		),
	)

	// Using AuthHandler wrapper
	mux.HandleFunc("/api/profile", middleware.AuthHandler(profileHandler))

	// Admin endpoint with custom logic
	mux.Handle("/api/admin", middleware.Authenticate(http.HandlerFunc(adminHandler)))

	// Health check
	mux.HandleFunc("/health", healthHandler)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /health              - Health check")
	fmt.Println("  GET  /api/public          - Public endpoint")
	fmt.Println("  GET  /api/protected       - Protected endpoint (requires auth)")
	fmt.Println("  GET  /api/optional        - Optional auth endpoint")
	fmt.Println("  GET  /api/verified-only   - Requires verified email")
	fmt.Println("  GET  /api/profile         - User profile (requires auth)")
	fmt.Println("  GET  /api/admin           - Admin endpoint (requires auth)")
	fmt.Println()
	fmt.Println("To test, use:")
	fmt.Println("  curl http://localhost:8080/api/public")
	fmt.Println("  curl -H 'Authorization: Bearer YOUR_TOKEN' http://localhost:8080/api/protected")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

// publicHandler handles public requests (no authentication required)
func publicHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "This is a public endpoint",
		"time":    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// protectedHandler handles authenticated requests
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	user := betterauth.GetUserFromContext(r.Context())
	session := betterauth.GetSessionFromContext(r.Context())

	if user == nil {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "This is a protected endpoint",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"session": map[string]interface{}{
			"id":        session.ID,
			"expiresAt": session.ExpiresAt.Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// optionalAuthHandler works with or without authentication
func optionalAuthHandler(w http.ResponseWriter, r *http.Request) {
	user := betterauth.GetUserFromContext(r.Context())

	response := map[string]interface{}{
		"message": "This endpoint works with or without authentication",
	}

	if user != nil {
		response["authenticated"] = true
		response["user"] = map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		}
	} else {
		response["authenticated"] = false
		response["user"] = nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// verifiedOnlyHandler requires a verified email address
func verifiedOnlyHandler(w http.ResponseWriter, r *http.Request) {
	user := betterauth.GetUserFromContext(r.Context())

	response := map[string]interface{}{
		"message":       "This endpoint requires a verified email",
		"emailVerified": user.EmailVerified,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// profileHandler uses the AuthHandler wrapper
func profileHandler(w http.ResponseWriter, r *http.Request, user *betterauth.User) {
	session := betterauth.GetSessionFromContext(r.Context())

	response := map[string]interface{}{
		"profile": map[string]interface{}{
			"id":            user.ID,
			"email":         user.Email,
			"name":          user.Name,
			"image":         user.Image,
			"emailVerified": user.EmailVerified,
			"createdAt":     user.CreatedAt.Format(time.RFC3339),
			"updatedAt":     user.UpdatedAt.Format(time.RFC3339),
		},
		"session": map[string]interface{}{
			"id":        session.ID,
			"expiresAt": session.ExpiresAt.Format(time.RFC3339),
			"ipAddress": session.IPAddress,
			"userAgent": session.UserAgent,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// adminHandler demonstrates custom authorization logic
func adminHandler(w http.ResponseWriter, r *http.Request) {
	user := betterauth.GetUserFromContext(r.Context())

	// Custom admin check - in a real app, you'd check user roles/permissions
	// For this example, we'll just check if metadata contains isAdmin flag
	isAdmin := false
	if user.Metadata != nil {
		if admin, ok := user.Metadata["isAdmin"].(bool); ok {
			isAdmin = admin
		}
	}

	if !isAdmin {
		http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	response := map[string]interface{}{
		"message": "Welcome to the admin panel",
		"admin": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// healthHandler returns the health status of the API
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
