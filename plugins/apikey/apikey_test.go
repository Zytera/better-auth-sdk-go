package apikey_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/apikey"
)

func TestCreate(t *testing.T) {
	var gotMethod, gotPath string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &body)
		w.Write([]byte(`{"id":"k1","name":"ci","key":"ba_secret","prefix":"ba_"}`))
	}))
	defer srv.Close()

	c := betterauth.NewClient(&betterauth.Config{BaseURL: srv.URL}, &betterauth.SessionToken{})
	p := apikey.New(c)

	key, err := p.Create(context.Background(), apikey.CreateInput{Name: "ci", Prefix: "ba_"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/api/auth/api-key/create" {
		t.Fatalf("bad route: %s %s", gotMethod, gotPath)
	}
	if body["name"] != "ci" || body["prefix"] != "ba_" {
		t.Fatalf("bad body: %v", body)
	}
	if key.ID != "k1" || key.Key != "ba_secret" {
		t.Fatalf("bad decode: %+v", key)
	}
}

func TestList(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotQuery = r.Method, r.URL.Path, r.URL.RawQuery
		w.Write([]byte(`{"apiKeys":[{"id":"k1"}],"total":1,"limit":10,"offset":5}`))
	}))
	defer srv.Close()

	c := betterauth.NewClient(&betterauth.Config{BaseURL: srv.URL}, &betterauth.SessionToken{})
	p := apikey.New(c)

	res, err := p.List(context.Background(), apikey.ListQuery{Limit: 10, Offset: 5, SortBy: "createdAt", SortDirection: "desc"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/api/auth/api-key/list" {
		t.Fatalf("bad route: %s %s", gotMethod, gotPath)
	}
	if gotQuery != "limit=10&offset=5&sortBy=createdAt&sortDirection=desc" {
		t.Fatalf("bad query: %s", gotQuery)
	}
	if res.Total != 1 || len(res.APIKeys) != 1 {
		t.Fatalf("bad decode: %+v", res)
	}
}

func TestGetValidation(t *testing.T) {
	c := betterauth.NewClient(&betterauth.Config{BaseURL: "http://localhost"}, &betterauth.SessionToken{})
	p := apikey.New(c)

	if _, err := p.Get(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty id")
	} else if !betterauth.IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}
