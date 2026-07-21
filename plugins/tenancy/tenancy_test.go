package tenancy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	betterauth "github.com/Zytera/better-auth-sdk-go"
	"github.com/Zytera/better-auth-sdk-go/plugins/tenancy"
)

func TestRoutesAndQuery(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotQuery = r.Method, r.URL.Path, r.URL.RawQuery
		switch r.URL.Path {
		case "/custom/auth/tenancy/organization/list":
			w.Write([]byte(`{"items":[{"id":"o1","name":"Acme"}],"total":1}`))
		case "/custom/auth/tenancy/member/check":
			w.Write([]byte(`{"isMember":true}`))
		case "/custom/auth/tenancy/permission/contexts":
			w.Write([]byte(`{"contexts":[{"contextType":"organization","contextId":"o1","name":"Acme","parentType":null,"parentId":null,"organizationId":"o1","source":"membership"}]}`))
		default:
			w.Write([]byte(`{"id":"x"}`))
		}
	}))
	defer srv.Close()

	// Custom BasePath must be applied by the core client to every plugin route.
	c := betterauth.NewClient(&betterauth.Config{BaseURL: srv.URL, BasePath: "/custom/auth"}, &betterauth.SessionToken{})
	p := tenancy.New(c)
	ctx := context.Background()

	// GET list with pagination query
	list, err := p.Organization.List(ctx, tenancy.ListQuery{Limit: 10, Offset: 5, OrderBy: "name"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/custom/auth/tenancy/organization/list" {
		t.Fatalf("bad route: %s %s", gotMethod, gotPath)
	}
	if gotQuery != "limit=10&offset=5&orderBy=name" {
		t.Fatalf("bad query: %q", gotQuery)
	}
	if list.Total != 1 || len(list.Items) != 1 || list.Items[0].Name != "Acme" {
		t.Fatalf("bad decode: %+v", list)
	}

	// GET by id lands on the right route
	if _, err := p.Team.Get(ctx, "t1"); err != nil {
		t.Fatalf("Team.Get: %v", err)
	}
	if gotPath != "/custom/auth/tenancy/team/t1" {
		t.Fatalf("bad route: %s", gotPath)
	}

	// POST returning a scalar field
	ok, err := p.Member.Check(ctx, "u1", tenancy.ContextOrganization, "o1")
	if err != nil || !ok {
		t.Fatalf("Check: ok=%v err=%v", ok, err)
	}
	if gotMethod != "POST" {
		t.Fatalf("Check should POST, got %s", gotMethod)
	}

	// GET permission/contexts: optional-param encoding + nullable-field decode.
	ctxs, err := p.Permission.Contexts(ctx, tenancy.ContextsQuery{
		StatementID: "doc:read", ContextType: tenancy.ContextOrganization, IncludeDescendants: true,
	})
	if err != nil {
		t.Fatalf("Contexts: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/custom/auth/tenancy/permission/contexts" {
		t.Fatalf("bad route: %s %s", gotMethod, gotPath)
	}
	if gotQuery != "contextType=organization&includeDescendants=true&statementId=doc%3Aread" {
		t.Fatalf("bad query: %q", gotQuery)
	}
	if len(ctxs) != 1 || ctxs[0].OrganizationID != "o1" || ctxs[0].ParentType != nil || ctxs[0].RoleID != nil {
		t.Fatalf("bad decode: %+v", ctxs)
	}
}
