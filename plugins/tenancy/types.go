package tenancy

import (
	"time"

	betterauth "github.com/Zytera/better-auth-sdk-go"
)

// ContextType is the hierarchy level a role/permission/member/invitation targets.
type ContextType = string

const (
	ContextOrganization ContextType = "organization"
	ContextTeam         ContextType = "team"
)

// List is the shape returned by the list endpoints: { items, total, limit, offset }.
type List[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Organization is the root entity of the hierarchy.
type Organization struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// Team is a hierarchical group nested under an organization or another team
// (unlimited nesting via parentType/parentId).
type Team struct {
	ID         string                 `json:"id"`
	ParentType ContextType            `json:"parentType"`
	ParentID   string                 `json:"parentId"`
	Name       string                 `json:"name"`
	Slug       string                 `json:"slug"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}

// Statement is a single permission, formatted as "category:operation".
type Statement struct {
	ID        string    `json:"id"`
	Category  string    `json:"category"`
	Operation string    `json:"operation"`
	CreatedAt time.Time `json:"createdAt"`
}

// Role is a custom role scoped to an organization or team context.
type Role struct {
	ID          string      `json:"id"`
	ContextType ContextType `json:"contextType"`
	ContextID   string      `json:"contextId"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	// Statements is populated by list when includeStatements=true.
	Statements []string `json:"statements,omitempty"`
}

// Member links a user to a context, optionally with a role.
type Member struct {
	ID          string      `json:"id"`
	ContextType ContextType `json:"contextType"`
	ContextID   string      `json:"contextId"`
	UserID      string      `json:"userId"`
	RoleID      *string     `json:"roleId,omitempty"`
	// User is only populated by endpoints that return enriched members
	// (they omit userId); nil elsewhere.
	User      *betterauth.User       `json:"user,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// PermissionGrant is an explicit permission granted to a user.
type PermissionGrant struct {
	ID          string      `json:"id"`
	UserID      string      `json:"userId"`
	StatementID string      `json:"statementId"`
	ContextType ContextType `json:"contextType"`
	ContextID   string      `json:"contextId"`
	ExpiresAt   *time.Time  `json:"expiresAt,omitempty"`
	GrantedBy   string      `json:"grantedBy"`
	CreatedAt   time.Time   `json:"createdAt"`
}

// PermissionDeny is an explicit permission denial (overrides everything).
type PermissionDeny struct {
	ID          string      `json:"id"`
	UserID      string      `json:"userId"`
	StatementID string      `json:"statementId"`
	ContextType ContextType `json:"contextType"`
	ContextID   string      `json:"contextId"`
	DeniedBy    string      `json:"deniedBy"`
	CreatedAt   time.Time   `json:"createdAt"`
}

// CheckResult is returned by permission/check.
type CheckResult struct {
	Allowed bool                   `json:"allowed"`
	Reason  string                 `json:"reason,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// CheckInput is the payload for Permission.Check.
type CheckInput struct {
	UserID      string      `json:"userId"`
	StatementID string      `json:"statementId"`
	ContextType ContextType `json:"contextType"`
	ContextID   string      `json:"contextId"`
}

// AccessibleContext is one entry from permission/contexts: a context the
// authenticated user can access, with its hierarchy and access source.
type AccessibleContext struct {
	ContextType    ContextType  `json:"contextType"`
	ContextID      string       `json:"contextId"`
	Name           string       `json:"name"`
	ParentType     *ContextType `json:"parentType,omitempty"`
	ParentID       *string      `json:"parentId,omitempty"`
	OrganizationID string       `json:"organizationId"`
	Source         string       `json:"source"` // "role" | "membership" | "grant" | "inherited"
	RoleID         *string      `json:"roleId,omitempty"`
}

// Invitation is a context invitation (email/phone/code/QR agnostic).
type Invitation struct {
	ID             string                 `json:"id"`
	Token          string                 `json:"token"`
	InvitationCode *string                `json:"invitationCode,omitempty"`
	Identifier     *string                `json:"identifier,omitempty"`
	ContextType    ContextType            `json:"contextType"`
	ContextID      string                 `json:"contextId"`
	RoleID         *string                `json:"roleId,omitempty"`
	InvitedBy      string                 `json:"invitedBy"`
	Status         string                 `json:"status"`
	ExpiresAt      *time.Time             `json:"expiresAt,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	AcceptedAt     *time.Time             `json:"acceptedAt,omitempty"`
	CancelledAt    *time.Time             `json:"cancelledAt,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
