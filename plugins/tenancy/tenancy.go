// Package tenancy is the client-side plugin for the @zytera/better-auth-tenancy
// multi-tenancy plugin: hierarchical organizations/teams, custom roles,
// statements, memberships, explicit permissions and invitations.
//
// Endpoints and payloads mirror the plugin's official TypeScript client
// (src/clients/index.ts). All routes live under /tenancy.
package tenancy

import (
	"context"
	"net/url"
	"strconv"
	"time"

	betterauth "github.com/Zytera/better-auth-sdk-go"
)

const routePrefix = "/tenancy"

// Plugin groups the tenancy endpoints into sub-services, matching the TS
// client (tenancy.organization.create, tenancy.team.list, ...).
type Plugin struct {
	Organization *organizationService
	Team         *teamService
	Statement    *statementService
	Role         *roleService
	Member       *memberService
	Permission   *permissionService
	Invitation   *invitationService
}

// New wires the plugin to any Requester (typically *betterauth.Client).
func New(r betterauth.Requester) *Plugin {
	return &Plugin{
		Organization: &organizationService{r},
		Team:         &teamService{r},
		Statement:    &statementService{r},
		Role:         &roleService{r},
		Member:       &memberService{r},
		Permission:   &permissionService{r},
		Invitation:   &invitationService{r},
	}
}

// --- validation helpers ------------------------------------------------------

func validationError(field string) error {
	return betterauth.NewError(betterauth.ErrorTypeValidation, field+" is required")
}

func requireString(field, value string) error {
	if value == "" {
		return validationError(field)
	}
	return nil
}

// --- query helpers -----------------------------------------------------------

func withQuery(path string, q url.Values) string {
	if len(q) == 0 {
		return path
	}
	return path + "?" + q.Encode()
}

func (q ListQuery) values() url.Values {
	v := url.Values{}
	if q.Limit > 0 {
		v.Set("limit", strconv.Itoa(q.Limit))
	}
	if q.Offset > 0 {
		v.Set("offset", strconv.Itoa(q.Offset))
	}
	if q.OrderBy != "" {
		v.Set("orderBy", q.OrderBy)
	}
	if q.OrderDirection != "" {
		v.Set("orderDirection", q.OrderDirection)
	}
	return v
}

// ListQuery holds the pagination/ordering options shared by list endpoints.
type ListQuery struct {
	Limit          int
	Offset         int
	OrderBy        string // "createdAt" | "name" | "id" | "addedAt" | "expiresAt" (per resource)
	OrderDirection string // "asc" | "desc"
}

// --- organization ------------------------------------------------------------

type organizationService struct{ r betterauth.Requester }

// CreateOrgInput is the payload for Organization.Create.
type CreateOrgInput struct {
	Name     string                 `json:"name"`
	Slug     string                 `json:"slug"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (s *organizationService) Create(ctx context.Context, in CreateOrgInput) (*Organization, error) {
	if in.Name == "" {
		return nil, validationError("name")
	}
	if in.Slug == "" {
		return nil, validationError("slug")
	}
	return do[Organization](s.r, ctx, "POST", routePrefix+"/organization/create", in)
}

func (s *organizationService) Get(ctx context.Context, id string) (*Organization, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	return do[Organization](s.r, ctx, "GET", routePrefix+"/organization/"+url.PathEscape(id), nil)
}

func (s *organizationService) GetBySlug(ctx context.Context, slug string) (*Organization, error) {
	if err := requireString("slug", slug); err != nil {
		return nil, err
	}
	return do[Organization](s.r, ctx, "GET", routePrefix+"/organization/by-slug/"+url.PathEscape(slug), nil)
}

// UpdateOrgInput is the payload for Organization.Update (ID required).
type UpdateOrgInput struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name,omitempty"`
	Slug     string                 `json:"slug,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (s *organizationService) Update(ctx context.Context, in UpdateOrgInput) (*Organization, error) {
	if err := requireString("id", in.ID); err != nil {
		return nil, err
	}
	return do[Organization](s.r, ctx, "POST", routePrefix+"/organization/update", in)
}

func (s *organizationService) Delete(ctx context.Context, id string) error {
	if err := requireString("id", id); err != nil {
		return err
	}
	return s.r.Do(ctx, "POST", routePrefix+"/organization/delete", map[string]string{"id": id}, nil)
}

func (s *organizationService) List(ctx context.Context, q ListQuery) (*List[Organization], error) {
	return do[List[Organization]](s.r, ctx, "GET", withQuery(routePrefix+"/organization/list", q.values()), nil)
}

// --- team --------------------------------------------------------------------

type teamService struct{ r betterauth.Requester }

// CreateTeamInput is the payload for Team.Create.
type CreateTeamInput struct {
	ParentType ContextType            `json:"parentType"`
	ParentID   string                 `json:"parentId"`
	Name       string                 `json:"name"`
	Slug       string                 `json:"slug"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

func (s *teamService) Create(ctx context.Context, in CreateTeamInput) (*Team, error) {
	if err := requireString("parentType", in.ParentType); err != nil {
		return nil, err
	}
	if err := requireString("parentId", in.ParentID); err != nil {
		return nil, err
	}
	if in.Name == "" {
		return nil, validationError("name")
	}
	if in.Slug == "" {
		return nil, validationError("slug")
	}
	return do[Team](s.r, ctx, "POST", routePrefix+"/team/create", in)
}

func (s *teamService) Get(ctx context.Context, id string) (*Team, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	return do[Team](s.r, ctx, "GET", routePrefix+"/team/"+url.PathEscape(id), nil)
}

func (s *teamService) GetBySlug(ctx context.Context, parentType ContextType, parentID, slug string) (*Team, error) {
	if err := requireString("parentType", parentType); err != nil {
		return nil, err
	}
	if err := requireString("parentId", parentID); err != nil {
		return nil, err
	}
	if err := requireString("slug", slug); err != nil {
		return nil, err
	}
	path := routePrefix + "/team/" + url.PathEscape(parentType) + "/" + url.PathEscape(parentID) + "/by-slug/" + url.PathEscape(slug)
	return do[Team](s.r, ctx, "GET", path, nil)
}

// UpdateTeamInput is the payload for Team.Update (ID required).
type UpdateTeamInput struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name,omitempty"`
	Slug     string                 `json:"slug,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (s *teamService) Update(ctx context.Context, in UpdateTeamInput) (*Team, error) {
	if err := requireString("id", in.ID); err != nil {
		return nil, err
	}
	return do[Team](s.r, ctx, "POST", routePrefix+"/team/update", in)
}

func (s *teamService) Delete(ctx context.Context, id string) error {
	if err := requireString("id", id); err != nil {
		return err
	}
	return s.r.Do(ctx, "POST", routePrefix+"/team/delete", map[string]string{"id": id}, nil)
}

func (s *teamService) List(ctx context.Context, parentType ContextType, parentID string, q ListQuery) (*List[Team], error) {
	if err := requireString("parentType", parentType); err != nil {
		return nil, err
	}
	if err := requireString("parentId", parentID); err != nil {
		return nil, err
	}
	v := q.values()
	v.Set("parentType", parentType)
	v.Set("parentId", parentID)
	return do[List[Team]](s.r, ctx, "GET", withQuery(routePrefix+"/team/list", v), nil)
}

// --- statement ---------------------------------------------------------------

type statementService struct{ r betterauth.Requester }

// Create registers a statement. id is "category:operation".
func (s *statementService) Create(ctx context.Context, id string) (*Statement, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	return do[Statement](s.r, ctx, "POST", routePrefix+"/statement/create", map[string]string{"id": id})
}

func (s *statementService) Get(ctx context.Context, id string) (*Statement, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	return do[Statement](s.r, ctx, "GET", routePrefix+"/statement/"+url.PathEscape(id), nil)
}

// UpdateStatementInput is the payload for Statement.Update.
type UpdateStatementInput struct {
	ID        string `json:"id"`
	Category  string `json:"category,omitempty"`
	Operation string `json:"operation,omitempty"`
}

func (s *statementService) Update(ctx context.Context, in UpdateStatementInput) (*Statement, error) {
	if err := requireString("id", in.ID); err != nil {
		return nil, err
	}
	return do[Statement](s.r, ctx, "POST", routePrefix+"/statement/update", in)
}

func (s *statementService) Delete(ctx context.Context, id string) error {
	if err := requireString("id", id); err != nil {
		return err
	}
	return s.r.Do(ctx, "POST", routePrefix+"/statement/delete", map[string]string{"id": id}, nil)
}

func (s *statementService) List(ctx context.Context, category string, q ListQuery) (*List[Statement], error) {
	v := q.values()
	if category != "" {
		v.Set("category", category)
	}
	return do[List[Statement]](s.r, ctx, "GET", withQuery(routePrefix+"/statement/list", v), nil)
}

func (s *statementService) BatchCreate(ctx context.Context, statements []string) ([]Statement, error) {
	if len(statements) == 0 {
		return nil, validationError("statements")
	}
	var out []Statement
	err := s.r.Do(ctx, "POST", routePrefix+"/statement/batch-create", map[string]interface{}{"statements": statements}, &out)
	return out, err
}

// --- role --------------------------------------------------------------------

type roleService struct{ r betterauth.Requester }

// CreateRoleInput is the payload for Role.Create.
type CreateRoleInput struct {
	Name        string      `json:"name"`
	ContextType ContextType `json:"contextType"`
	ContextID   string      `json:"contextId"`
	Statements  []string    `json:"statements,omitempty"`
	Description string      `json:"description,omitempty"`
}

func (s *roleService) Create(ctx context.Context, in CreateRoleInput) (*Role, error) {
	if in.Name == "" {
		return nil, validationError("name")
	}
	if err := requireString("contextType", in.ContextType); err != nil {
		return nil, err
	}
	if err := requireString("contextId", in.ContextID); err != nil {
		return nil, err
	}
	return do[Role](s.r, ctx, "POST", routePrefix+"/role/create", in)
}

func (s *roleService) Get(ctx context.Context, id string) (*Role, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	return do[Role](s.r, ctx, "GET", routePrefix+"/role/"+url.PathEscape(id), nil)
}

// UpdateRoleInput is the payload for Role.Update.
type UpdateRoleInput struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (s *roleService) Update(ctx context.Context, in UpdateRoleInput) (*Role, error) {
	if err := requireString("id", in.ID); err != nil {
		return nil, err
	}
	return do[Role](s.r, ctx, "POST", routePrefix+"/role/update", in)
}

func (s *roleService) Delete(ctx context.Context, id string) error {
	if err := requireString("id", id); err != nil {
		return err
	}
	return s.r.Do(ctx, "POST", routePrefix+"/role/delete", map[string]string{"id": id}, nil)
}

func (s *roleService) List(ctx context.Context, contextType ContextType, contextID string, includeStatements bool, q ListQuery) (*List[Role], error) {
	if err := requireString("contextType", contextType); err != nil {
		return nil, err
	}
	if err := requireString("contextId", contextID); err != nil {
		return nil, err
	}
	v := q.values()
	v.Set("contextType", contextType)
	v.Set("contextId", contextID)
	if includeStatements {
		v.Set("includeStatements", "true")
	}
	return do[List[Role]](s.r, ctx, "GET", withQuery(routePrefix+"/role/list", v), nil)
}

func (s *roleService) AddStatement(ctx context.Context, roleID, statementID string) (*Role, error) {
	if err := requireString("roleID", roleID); err != nil {
		return nil, err
	}
	if err := requireString("statementID", statementID); err != nil {
		return nil, err
	}
	return do[Role](s.r, ctx, "POST", routePrefix+"/role/statement/add", map[string]string{
		"id": roleID, "statementId": statementID,
	})
}

func (s *roleService) RemoveStatement(ctx context.Context, roleID, statementID string) (*Role, error) {
	if err := requireString("roleID", roleID); err != nil {
		return nil, err
	}
	if err := requireString("statementID", statementID); err != nil {
		return nil, err
	}
	return do[Role](s.r, ctx, "POST", routePrefix+"/role/statement/remove", map[string]string{
		"id": roleID, "statementId": statementID,
	})
}

// --- member ------------------------------------------------------------------

type memberService struct{ r betterauth.Requester }

// AddMemberInput is the payload for Member.Add.
type AddMemberInput struct {
	UserID      string                 `json:"userId"`
	ContextType ContextType            `json:"contextType"`
	ContextID   string                 `json:"contextId"`
	RoleID      string                 `json:"roleId,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func (s *memberService) Add(ctx context.Context, in AddMemberInput) (*Member, error) {
	if err := requireString("userId", in.UserID); err != nil {
		return nil, err
	}
	if err := requireString("contextType", in.ContextType); err != nil {
		return nil, err
	}
	if err := requireString("contextId", in.ContextID); err != nil {
		return nil, err
	}
	return do[Member](s.r, ctx, "POST", routePrefix+"/member/add", in)
}

func (s *memberService) Get(ctx context.Context, id string) (*Member, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	return do[Member](s.r, ctx, "GET", routePrefix+"/member/"+url.PathEscape(id), nil)
}

func (s *memberService) UpdateRole(ctx context.Context, id, roleID, updatedBy string) (*Member, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	if err := requireString("roleID", roleID); err != nil {
		return nil, err
	}
	if err := requireString("updatedBy", updatedBy); err != nil {
		return nil, err
	}
	return do[Member](s.r, ctx, "POST", routePrefix+"/member/role/update", map[string]string{
		"id": id, "roleId": roleID, "updatedBy": updatedBy,
	})
}

func (s *memberService) UpdateMetadata(ctx context.Context, id string, metadata map[string]interface{}, updatedBy string) (*Member, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	if err := requireString("updatedBy", updatedBy); err != nil {
		return nil, err
	}
	return do[Member](s.r, ctx, "POST", routePrefix+"/member/metadata/update", map[string]interface{}{
		"id": id, "metadata": metadata, "updatedBy": updatedBy,
	})
}

func (s *memberService) Remove(ctx context.Context, id, removedBy string) error {
	if err := requireString("id", id); err != nil {
		return err
	}
	if err := requireString("removedBy", removedBy); err != nil {
		return err
	}
	return s.r.Do(ctx, "POST", routePrefix+"/member/remove", map[string]string{
		"id": id, "removedBy": removedBy,
	}, nil)
}

func (s *memberService) List(ctx context.Context, contextType ContextType, contextID string, q ListQuery) (*List[Member], error) {
	if err := requireString("contextType", contextType); err != nil {
		return nil, err
	}
	if err := requireString("contextId", contextID); err != nil {
		return nil, err
	}
	v := q.values()
	v.Set("contextType", contextType)
	v.Set("contextId", contextID)
	return do[List[Member]](s.r, ctx, "GET", withQuery(routePrefix+"/member/list", v), nil)
}

// GetUserMemberships lists a user's memberships. contextType is optional ("").
func (s *memberService) GetUserMemberships(ctx context.Context, userID string, contextType ContextType, limit, offset int) (*List[Member], error) {
	if err := requireString("userId", userID); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("userId", userID)
	if contextType != "" {
		v.Set("contextType", contextType)
	}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		v.Set("offset", strconv.Itoa(offset))
	}
	return do[List[Member]](s.r, ctx, "GET", withQuery(routePrefix+"/member/user/memberships", v), nil)
}

// Check reports whether the user is a member of the given context.
func (s *memberService) Check(ctx context.Context, userID string, contextType ContextType, contextID string) (bool, error) {
	if err := requireString("userId", userID); err != nil {
		return false, err
	}
	if err := requireString("contextType", contextType); err != nil {
		return false, err
	}
	if err := requireString("contextId", contextID); err != nil {
		return false, err
	}
	var out struct {
		IsMember bool `json:"isMember"`
	}
	err := s.r.Do(ctx, "POST", routePrefix+"/member/check", map[string]string{
		"userId": userID, "contextType": contextType, "contextId": contextID,
	}, &out)
	return out.IsMember, err
}

// --- permission --------------------------------------------------------------

type permissionService struct{ r betterauth.Requester }

// GrantInput is the payload for Permission.Grant.
type GrantInput struct {
	UserID      string      `json:"userId"`
	StatementID string      `json:"statementId"`
	ContextType ContextType `json:"contextType"`
	ContextID   string      `json:"contextId"`
	ExpiresAt   *time.Time  `json:"expiresAt,omitempty"`
}

func (s *permissionService) Grant(ctx context.Context, in GrantInput) (*PermissionGrant, error) {
	if err := requireString("userId", in.UserID); err != nil {
		return nil, err
	}
	if err := requireString("statementId", in.StatementID); err != nil {
		return nil, err
	}
	if err := requireString("contextType", in.ContextType); err != nil {
		return nil, err
	}
	if err := requireString("contextId", in.ContextID); err != nil {
		return nil, err
	}
	return do[PermissionGrant](s.r, ctx, "POST", routePrefix+"/permission/grant", in)
}

// DenyInput is the payload for Permission.Deny.
type DenyInput struct {
	UserID      string      `json:"userId"`
	StatementID string      `json:"statementId"`
	ContextType ContextType `json:"contextType"`
	ContextID   string      `json:"contextId"`
}

func (s *permissionService) Deny(ctx context.Context, in DenyInput) (*PermissionDeny, error) {
	if err := requireString("userId", in.UserID); err != nil {
		return nil, err
	}
	if err := requireString("statementId", in.StatementID); err != nil {
		return nil, err
	}
	if err := requireString("contextType", in.ContextType); err != nil {
		return nil, err
	}
	if err := requireString("contextId", in.ContextID); err != nil {
		return nil, err
	}
	return do[PermissionDeny](s.r, ctx, "POST", routePrefix+"/permission/deny", in)
}

// Revoke removes a grant or deny. typ is "grant" or "deny".
func (s *permissionService) Revoke(ctx context.Context, permissionID, typ string) error {
	if err := requireString("permissionID", permissionID); err != nil {
		return err
	}
	if err := requireString("type", typ); err != nil {
		return err
	}
	return s.r.Do(ctx, "POST", routePrefix+"/permission/revoke", map[string]string{
		"permissionId": permissionID, "type": typ,
	}, nil)
}

// Check evaluates whether a user has a statement in a context.
func (s *permissionService) Check(ctx context.Context, in CheckInput) (*CheckResult, error) {
	if err := requireString("userId", in.UserID); err != nil {
		return nil, err
	}
	if err := requireString("statementId", in.StatementID); err != nil {
		return nil, err
	}
	if err := requireString("contextType", in.ContextType); err != nil {
		return nil, err
	}
	if err := requireString("contextId", in.ContextID); err != nil {
		return nil, err
	}
	return do[CheckResult](s.r, ctx, "POST", routePrefix+"/permission/check", in)
}

func (s *permissionService) ListGrants(ctx context.Context, userID string, limit, offset int) ([]PermissionGrant, error) {
	if err := requireString("userId", userID); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("userId", userID)
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		v.Set("offset", strconv.Itoa(offset))
	}
	var out []PermissionGrant
	err := s.r.Do(ctx, "GET", withQuery(routePrefix+"/permission/grants", v), nil, &out)
	return out, err
}

func (s *permissionService) ListDenies(ctx context.Context, userID string, limit, offset int) ([]PermissionDeny, error) {
	if err := requireString("userId", userID); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("userId", userID)
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		v.Set("offset", strconv.Itoa(offset))
	}
	var out []PermissionDeny
	err := s.r.Do(ctx, "GET", withQuery(routePrefix+"/permission/denies", v), nil, &out)
	return out, err
}

func (s *permissionService) GrantBatch(ctx context.Context, grants []GrantInput) ([]PermissionGrant, error) {
	if len(grants) == 0 {
		return nil, validationError("grants")
	}
	var out []PermissionGrant
	err := s.r.Do(ctx, "POST", routePrefix+"/permission/grant-batch", map[string]interface{}{"grants": grants}, &out)
	return out, err
}

func (s *permissionService) DenyBatch(ctx context.Context, denies []DenyInput) ([]PermissionDeny, error) {
	if len(denies) == 0 {
		return nil, validationError("denies")
	}
	var out []PermissionDeny
	err := s.r.Do(ctx, "POST", routePrefix+"/permission/deny-batch", map[string]interface{}{"denies": denies}, &out)
	return out, err
}

// --- invitation --------------------------------------------------------------

type invitationService struct{ r betterauth.Requester }

// CreateInvitationInput is the payload for Invitation.Create.
type CreateInvitationInput struct {
	Identifier     string                 `json:"identifier,omitempty"`
	ContextType    ContextType            `json:"contextType"`
	ContextID      string                 `json:"contextId"`
	RoleID         string                 `json:"roleId,omitempty"`
	ExpiresInHours int                    `json:"expiresInHours,omitempty"`
	GenerateCode   bool                   `json:"generateCode,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

func (s *invitationService) Create(ctx context.Context, in CreateInvitationInput) (*Invitation, error) {
	if err := requireString("contextType", in.ContextType); err != nil {
		return nil, err
	}
	if err := requireString("contextId", in.ContextID); err != nil {
		return nil, err
	}
	return do[Invitation](s.r, ctx, "POST", routePrefix+"/invitation/create", in)
}

func (s *invitationService) Accept(ctx context.Context, token string) (*Invitation, error) {
	if err := requireString("token", token); err != nil {
		return nil, err
	}
	return do[Invitation](s.r, ctx, "POST", routePrefix+"/invitation/accept", map[string]string{"token": token})
}

func (s *invitationService) Reject(ctx context.Context, token, userID string) (*Invitation, error) {
	if err := requireString("token", token); err != nil {
		return nil, err
	}
	if err := requireString("userId", userID); err != nil {
		return nil, err
	}
	return do[Invitation](s.r, ctx, "POST", routePrefix+"/invitation/reject", map[string]string{
		"token": token, "userId": userID,
	})
}

func (s *invitationService) Cancel(ctx context.Context, id string) (*Invitation, error) {
	if err := requireString("id", id); err != nil {
		return nil, err
	}
	return do[Invitation](s.r, ctx, "POST", routePrefix+"/invitation/cancel", map[string]string{"id": id})
}

func (s *invitationService) GetByToken(ctx context.Context, token string) (*Invitation, error) {
	if err := requireString("token", token); err != nil {
		return nil, err
	}
	return do[Invitation](s.r, ctx, "GET", routePrefix+"/invitation/by-token/"+url.PathEscape(token), nil)
}

func (s *invitationService) GetByCode(ctx context.Context, code string) (*Invitation, error) {
	if err := requireString("code", code); err != nil {
		return nil, err
	}
	return do[Invitation](s.r, ctx, "GET", routePrefix+"/invitation/by-code/"+url.PathEscape(code), nil)
}

func (s *invitationService) List(ctx context.Context, contextType ContextType, contextID, status string, q ListQuery) (*List[Invitation], error) {
	if err := requireString("contextType", contextType); err != nil {
		return nil, err
	}
	if err := requireString("contextId", contextID); err != nil {
		return nil, err
	}
	v := q.values()
	v.Set("contextType", contextType)
	v.Set("contextId", contextID)
	if status != "" {
		v.Set("status", status)
	}
	return do[List[Invitation]](s.r, ctx, "GET", withQuery(routePrefix+"/invitation/list", v), nil)
}

func (s *invitationService) GetByIdentifier(ctx context.Context, identifier, status string) (*List[Invitation], error) {
	if err := requireString("identifier", identifier); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("identifier", identifier)
	if status != "" {
		v.Set("status", status)
	}
	return do[List[Invitation]](s.r, ctx, "GET", withQuery(routePrefix+"/invitation/by-identifier", v), nil)
}

// CleanupExpired marks expired invitations. Server route is
// /tenancy/invitation/cleanup-expired (the TS client's /cleanup path is a
// known mismatch; this uses the route the server actually registers).
func (s *invitationService) CleanupExpired(ctx context.Context) error {
	return s.r.Do(ctx, "POST", routePrefix+"/invitation/cleanup-expired", map[string]interface{}{}, nil)
}

// do performs the request and decodes the JSON body into *T, returning nil on error.
func do[T any](r betterauth.Requester, ctx context.Context, method, path string, body interface{}) (*T, error) {
	var out T
	if err := r.Do(ctx, method, path, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
