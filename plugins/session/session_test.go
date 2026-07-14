package session_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/session"
)

func TestGetSession(t *testing.T) {
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.Write([]byte(`{"user":{"id":"u1","email":"a@b.com"},"session":{"id":"s1","userId":"u1","token":"tok","expiresAt":"2026-01-01T00:00:00Z","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z"}}`))
	}))
	defer srv.Close()

	c := betterauth.NewClient(
		&betterauth.Config{BaseURL: srv.URL},
		&betterauth.SessionToken{},
	)
	p := session.New(c)

	data, err := p.Get(context.Background())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/api/auth/get-session" {
		t.Fatalf("bad route: %s %s", gotMethod, gotPath)
	}
	if data.User.ID != "u1" || data.User.Email != "a@b.com" {
		t.Fatalf("bad decode: %+v", data.User)
	}
	if data.Session.ID != "s1" {
		t.Fatalf("bad session decode: %+v", data.Session)
	}
}

func TestVerifySession(t *testing.T) {
	var gotMethod, gotPath string
	var body map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &body)
		w.Write([]byte(`{"id":"s1","userId":"u1","token":"tok","expiresAt":"2026-01-01T00:00:00Z","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	c := betterauth.NewClient(
		&betterauth.Config{BaseURL: srv.URL},
		&betterauth.SessionToken{},
	)
	p := session.New(c)

	sess, err := p.Verify(context.Background(), "my-token")
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/api/auth/session/verify" {
		t.Fatalf("bad route: %s %s", gotMethod, gotPath)
	}
	if body["token"] != "my-token" {
		t.Fatalf("bad token body: %v", body)
	}
	if sess.ID != "s1" {
		t.Fatalf("bad decode: %+v", sess)
	}
}

func TestVerifySessionRequiresToken(t *testing.T) {
	c := betterauth.NewClient(
		&betterauth.Config{BaseURL: "http://localhost"},
		&betterauth.SessionToken{},
	)
	p := session.New(c)

	_, err := p.Verify(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
	if !betterauth.IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestGetSessionWithCookie(t *testing.T) {
	var gotCookie *http.Cookie
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, c := range r.Cookies() {
			if c.Name == "better-auth.session_token" {
				gotCookie = c
			}
		}
		w.Write([]byte(`{"user":{"id":"u1"},"session":{"id":"s1","userId":"u1","token":"tok","expiresAt":"2026-01-01T00:00:00Z","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z"}}`))
	}))
	defer srv.Close()

	c := betterauth.NewClient(
		&betterauth.Config{BaseURL: srv.URL},
		&betterauth.SessionToken{Cookie: &http.Cookie{
			Name:  "better-auth.session_token",
			Value: "session-value",
		}},
	)
	p := session.New(c)

	if _, err := p.Get(context.Background()); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if gotCookie == nil || gotCookie.Value != "session-value" {
		t.Fatalf("cookie not forwarded: %+v", gotCookie)
	}
}
