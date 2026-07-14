// Package admin is the client-side plugin for the better-auth admin plugin.
package admin

import (
	"context"

	betterauth "github.com/Zytera/better-auth-sdk-go"
)

// Plugin talks to the admin endpoints. Construct it with New(client).
type Plugin struct {
	r betterauth.Requester
}

// New wires the plugin to any Requester (typically *betterauth.Client).
func New(r betterauth.Requester) *Plugin {
	return &Plugin{r: r}
}

// CreateUserInput is the payload for CreateUser.
type CreateUserInput struct {
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Name     string                 `json:"name"`
	Role     string                 `json:"role,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// CreateUser creates a user with an explicit role.
func (p *Plugin) CreateUser(ctx context.Context, in CreateUserInput) (*betterauth.User, error) {
	if in.Email == "" {
		return nil, betterauth.NewError(betterauth.ErrorTypeValidation, "email is required")
	}
	if in.Password == "" {
		return nil, betterauth.NewError(betterauth.ErrorTypeValidation, "password is required")
	}

	var out struct {
		User betterauth.User `json:"user"`
	}
	if err := p.r.Do(ctx, "POST", "/admin/create-user", in, &out); err != nil {
		return nil, err
	}
	return &out.User, nil
}

// ListUsersResult is the response of ListUsers.
type ListUsersResult struct {
	Users  []betterauth.User `json:"users"`
	Total  int               `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

// ListUsersQuery holds the supported filters for ListUsers.
type ListUsersQuery struct {
	Limit       int
	Offset      int
	SearchValue string
	SearchField string
}

// ListUsers lists users. Pass a zero-value query for defaults.
func (p *Plugin) ListUsers(ctx context.Context, query ListUsersQuery) (*ListUsersResult, error) {
	body := map[string]interface{}{}
	if query.Limit > 0 {
		body["limit"] = query.Limit
	}
	if query.Offset > 0 {
		body["offset"] = query.Offset
	}
	if query.SearchValue != "" {
		body["searchValue"] = query.SearchValue
	}
	if query.SearchField != "" {
		body["searchField"] = query.SearchField
	}

	var out ListUsersResult
	if err := p.r.Do(ctx, "POST", "/admin/list-users", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetRole sets a user's role.
func (p *Plugin) SetRole(ctx context.Context, userID, role string) error {
	if userID == "" {
		return betterauth.NewError(betterauth.ErrorTypeValidation, "userID is required")
	}
	if role == "" {
		return betterauth.NewError(betterauth.ErrorTypeValidation, "role is required")
	}
	return p.r.Do(ctx, "POST", "/admin/set-role", map[string]string{
		"userId": userID,
		"role":   role,
	}, nil)
}

// SetUserPassword sets a new password for a user.
func (p *Plugin) SetUserPassword(ctx context.Context, userID, newPassword string) error {
	if userID == "" {
		return betterauth.NewError(betterauth.ErrorTypeValidation, "userID is required")
	}
	if newPassword == "" {
		return betterauth.NewError(betterauth.ErrorTypeValidation, "newPassword is required")
	}
	return p.r.Do(ctx, "POST", "/admin/set-user-password", map[string]string{
		"userId":      userID,
		"newPassword": newPassword,
	}, nil)
}

// BanUser bans a user. reason and expiresIn (seconds) are optional (pass "" / 0).
func (p *Plugin) BanUser(ctx context.Context, userID, reason string, expiresIn int) error {
	if userID == "" {
		return betterauth.NewError(betterauth.ErrorTypeValidation, "userID is required")
	}
	body := map[string]interface{}{"userId": userID}
	if reason != "" {
		body["banReason"] = reason
	}
	if expiresIn > 0 {
		body["banExpiresIn"] = expiresIn
	}
	return p.r.Do(ctx, "POST", "/admin/ban-user", body, nil)
}

// UnbanUser lifts a ban.
func (p *Plugin) UnbanUser(ctx context.Context, userID string) error {
	if userID == "" {
		return betterauth.NewError(betterauth.ErrorTypeValidation, "userID is required")
	}
	return p.r.Do(ctx, "POST", "/admin/unban-user", map[string]string{
		"userId": userID,
	}, nil)
}

// RemoveUser hard-deletes a user.
func (p *Plugin) RemoveUser(ctx context.Context, userID string) error {
	if userID == "" {
		return betterauth.NewError(betterauth.ErrorTypeValidation, "userID is required")
	}
	return p.r.Do(ctx, "POST", "/admin/remove-user", map[string]string{
		"userId": userID,
	}, nil)
}

// ImpersonateUser starts an impersonation session for the given user.
func (p *Plugin) ImpersonateUser(ctx context.Context, userID string) (*betterauth.SessionData, error) {
	if userID == "" {
		return nil, betterauth.NewError(betterauth.ErrorTypeValidation, "userID is required")
	}
	var out betterauth.SessionData
	if err := p.r.Do(ctx, "POST", "/admin/impersonate-user", map[string]string{
		"userId": userID,
	}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// StopImpersonating ends the current impersonation session.
func (p *Plugin) StopImpersonating(ctx context.Context) error {
	return p.r.Do(ctx, "POST", "/admin/stop-impersonating", nil, nil)
}
