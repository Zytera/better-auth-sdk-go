package admin_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/admin"
)

func TestCreateUser(t *testing.T) {
	var gotMethod, gotPath string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &body)
		w.Write([]byte(`{"user":{"id":"u1","email":"user@example.com","role":"user"}}`))
	}))
	defer srv.Close()

	c := betterauth.NewClient(
		&betterauth.Config{BaseURL: srv.URL},
		&betterauth.SessionToken{},
	)
	p := admin.New(c)

	user, err := p.CreateUser(context.Background(), admin.CreateUserInput{
		Email:    "user@example.com",
		Password: "s3cret!",
		Name:     "New User",
		Role:     "user",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/api/auth/admin/create-user" {
		t.Fatalf("bad route: %s %s", gotMethod, gotPath)
	}
	if body["email"] != "user@example.com" || body["password"] != "s3cret!" {
		t.Fatalf("bad body: %v", body)
	}
	if user.ID != "u1" {
		t.Fatalf("bad decode: %+v", user)
	}
}

func TestCreateUserValidation(t *testing.T) {
	c := betterauth.NewClient(
		&betterauth.Config{BaseURL: "http://localhost"},
		&betterauth.SessionToken{},
	)
	p := admin.New(c)

	_, err := p.CreateUser(context.Background(), admin.CreateUserInput{Password: "pass"})
	if err == nil {
		t.Fatal("expected error for missing email")
	}
	if !betterauth.IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}

	_, err = p.CreateUser(context.Background(), admin.CreateUserInput{Email: "a@b.com"})
	if err == nil {
		t.Fatal("expected error for missing password")
	}
	if !betterauth.IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestSetRoleValidation(t *testing.T) {
	c := betterauth.NewClient(
		&betterauth.Config{BaseURL: "http://localhost"},
		&betterauth.SessionToken{},
	)
	p := admin.New(c)

	if err := p.SetRole(context.Background(), "", "admin"); err == nil {
		t.Fatal("expected error for empty userID")
	}
	if err := p.SetRole(context.Background(), "u1", ""); err == nil {
		t.Fatal("expected error for empty role")
	}
}

func TestListUsers(t *testing.T) {
	var gotMethod, gotPath string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &body)
		w.Write([]byte(`{"users":[{"id":"u1"}],"total":1,"limit":10,"offset":5}`))
	}))
	defer srv.Close()

	c := betterauth.NewClient(
		&betterauth.Config{BaseURL: srv.URL},
		&betterauth.SessionToken{},
	)
	p := admin.New(c)

	res, err := p.ListUsers(context.Background(), admin.ListUsersQuery{
		Limit:       10,
		Offset:      5,
		SearchValue: "john",
		SearchField: "email",
	})
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/api/auth/admin/list-users" {
		t.Fatalf("bad route: %s %s", gotMethod, gotPath)
	}
	if body["limit"].(float64) != 10 || body["searchValue"] != "john" {
		t.Fatalf("bad body: %v", body)
	}
	if res.Total != 1 || len(res.Users) != 1 {
		t.Fatalf("bad decode: %+v", res)
	}
}
