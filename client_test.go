package betterauth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	betterauth "github.com/Zytera/better-auth-sdk-go"
)

func TestSetAPIKeyHeader(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("x-api-key")
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := betterauth.NewClient(&betterauth.Config{BaseURL: srv.URL}, &betterauth.SessionToken{})

	// No key set: header absent.
	if err := c.Do(context.Background(), "GET", "/get-session", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotKey != "" {
		t.Fatalf("expected no x-api-key header, got %q", gotKey)
	}

	// Key set: header present.
	c.SetAPIKey("secret-key")
	if err := c.Do(context.Background(), "GET", "/get-session", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotKey != "secret-key" {
		t.Fatalf("expected x-api-key=secret-key, got %q", gotKey)
	}

	// Cleared: header absent again.
	c.SetAPIKey("")
	if err := c.Do(context.Background(), "GET", "/get-session", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotKey != "" {
		t.Fatalf("expected cleared x-api-key header, got %q", gotKey)
	}
}

func TestSetAPIKeyCustomHeader(t *testing.T) {
	var gotDefault, gotCustom string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotDefault = r.Header.Get("x-api-key")
		gotCustom = r.Header.Get("x-my-api-key")
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := betterauth.NewClient(
		&betterauth.Config{BaseURL: srv.URL, APIKeyHeader: "x-my-api-key"},
		&betterauth.SessionToken{},
	)
	c.SetAPIKey("secret-key")
	if err := c.Do(context.Background(), "GET", "/get-session", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotCustom != "secret-key" {
		t.Fatalf("expected x-my-api-key=secret-key, got %q", gotCustom)
	}
	if gotDefault != "" {
		t.Fatalf("expected no default x-api-key header, got %q", gotDefault)
	}
}
