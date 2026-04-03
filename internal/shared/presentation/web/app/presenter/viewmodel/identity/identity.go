package identityviewmodel

import (
	"strings"
	"time"

	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	webflash "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/flash"
)

type SessionItem struct {
	IDSession int64
	IDUser    int64
	UserEmail string
	IP        string
	UserAgent string
	CreatedAt string
	ExpiresAt string
	UpdatedAt string
	Status    string
	IsCurrent bool
}

type SessionsPage struct {
	ViewerID          int64
	CanManageSessions bool
	OwnSessions       []SessionItem
	AdminSessions     []SessionItem
	AdminFilter       string
	Error             string
	Flash             webflash.Message
}

type UserRole struct {
	ID         int16
	Name       string
	AssignedBy int64
}

type UserListItem struct {
	ID          int64
	FullName    string
	Email       string
	Status      string
	LastLoginAt string
	Roles       []UserRole
}

type RoleOption struct {
	ID   int16
	Name string
}

type UsersPage struct {
	Filter        string
	CanManageRBAC bool
	Roles         []RoleOption
	Flash         webflash.Message
	Items         []UserListItem
}

type PermissionListItem struct {
	ID          int16
	Resource    string
	Action      string
	Code        string
	Description string
	CreatedAt   string
}

type PermissionForm struct {
	ID          int16
	Resource    string
	Action      string
	Description string
}

type RolePermission struct {
	ID          int16
	Code        string
	Description string
}

type RolePermissionBlock struct {
	ID          int16
	Name        string
	Description string
	Permissions []RolePermission
}

type PermissionsPage struct {
	Filter        string
	CanManageRBAC bool
	Roles         []RolePermissionBlock
	Permissions   []PermissionListItem
	Form          PermissionForm
	Flash         webflash.Message
}

type RoleListItem struct {
	ID          int16
	Name        string
	Description string
	CreatedAt   string
}

type RoleForm struct {
	ID          int16
	Name        string
	Description string
}

type RolesPage struct {
	CanManageRBAC bool
	Roles         []RoleListItem
	Form          RoleForm
	Flash         webflash.Message
}

func BuildSessionsPage(viewerID int64, canManageSessions bool, adminFilter string, flash webflash.Message, own []identityread.SessionReadModel, admin []identityread.SessionReadModel, currentSessionID int64, now time.Time) SessionsPage {
	return SessionsPage{
		ViewerID:          viewerID,
		CanManageSessions: canManageSessions,
		OwnSessions:       makeSessionViews(own, currentSessionID, now),
		AdminSessions:     makeSessionViews(admin, currentSessionID, now),
		AdminFilter:       strings.TrimSpace(adminFilter),
		Flash:             flash,
	}
}

func BuildUsersPage(items []identityread.UserListItemReadModel, filter string, canManageRBAC bool, roles []identityread.RoleListItemReadModel, assignments map[int64][]identityread.UserRoleAssignmentReadModel, flash webflash.Message) UsersPage {
	data := UsersPage{
		Filter:        strings.TrimSpace(filter),
		CanManageRBAC: canManageRBAC,
		Flash:         flash,
		Items:         make([]UserListItem, 0, len(items)),
	}
	if canManageRBAC {
		data.Roles = makeRoleOptions(roles)
	}
	for _, item := range items {
		status := "Ativo"
		if item.IsBlocked {
			status = "Bloqueado"
		}
		lastLogin := "Nunca"
		if item.LastLoginAt != nil {
			lastLogin = item.LastLoginAt.Format(time.DateTime)
		}
		view := UserListItem{
			ID:          item.ID,
			FullName:    fallback(item.FullName, "-"),
			Email:       fallback(item.Email, "-"),
			Status:      status,
			LastLoginAt: lastLogin,
		}
		if canManageRBAC {
			view.Roles = makeUserRoleViews(assignments[item.ID])
		}
		data.Items = append(data.Items, view)
	}
	return data
}

func BuildPermissionsPage(items []identityread.PermissionListItemReadModel, filter string, canManageRBAC bool, roles []identityread.RoleListItemReadModel, assignments map[int16][]identityread.RolePermissionAssignmentReadModel, form PermissionForm, flash webflash.Message) PermissionsPage {
	data := PermissionsPage{
		Filter:        strings.TrimSpace(filter),
		CanManageRBAC: canManageRBAC,
		Permissions:   make([]PermissionListItem, 0, len(items)),
		Form:          form,
		Flash:         flash,
	}
	for _, item := range items {
		data.Permissions = append(data.Permissions, PermissionListItem{
			ID:          item.ID,
			Resource:    item.Resource,
			Action:      item.Action,
			Code:        item.Resource + ":" + item.Action,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format(time.DateTime),
		})
	}
	if !canManageRBAC {
		return data
	}
	data.Roles = make([]RolePermissionBlock, 0, len(roles))
	for _, role := range roles {
		data.Roles = append(data.Roles, RolePermissionBlock{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: makeRolePermissionViews(assignments[role.ID]),
		})
	}
	return data
}

func BuildRolesPage(items []identityread.RoleListItemReadModel, canManageRBAC bool, form RoleForm, flash webflash.Message) RolesPage {
	data := RolesPage{
		CanManageRBAC: canManageRBAC,
		Roles:         make([]RoleListItem, 0, len(items)),
		Form:          form,
		Flash:         flash,
	}
	for _, item := range items {
		data.Roles = append(data.Roles, RoleListItem{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format(time.DateTime),
		})
	}
	return data
}

func makeSessionViews(items []identityread.SessionReadModel, currentSessionID int64, now time.Time) []SessionItem {
	out := make([]SessionItem, 0, len(items))
	for _, item := range items {
		out = append(out, SessionItem{
			IDSession: item.IDSession,
			IDUser:    item.IDUser,
			UserEmail: item.UserEmail,
			IP:        fallback(item.IP, "-"),
			UserAgent: fallback(item.UserAgent, "-"),
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04"),
			ExpiresAt: item.ExpiresAt.Format("2006-01-02 15:04"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04"),
			Status:    sessionStatus(item, now),
			IsCurrent: item.IDSession == currentSessionID,
		})
	}
	return out
}

func makeRoleOptions(items []identityread.RoleListItemReadModel) []RoleOption {
	out := make([]RoleOption, 0, len(items))
	for _, item := range items {
		out = append(out, RoleOption{ID: item.ID, Name: item.Name})
	}
	return out
}

func makeUserRoleViews(items []identityread.UserRoleAssignmentReadModel) []UserRole {
	out := make([]UserRole, 0, len(items))
	for _, item := range items {
		out = append(out, UserRole{
			ID:         item.IDRole,
			Name:       item.RoleName,
			AssignedBy: item.AssignedBy,
		})
	}
	return out
}

func makeRolePermissionViews(items []identityread.RolePermissionAssignmentReadModel) []RolePermission {
	out := make([]RolePermission, 0, len(items))
	for _, item := range items {
		out = append(out, RolePermission{
			ID:          item.IDPermission,
			Code:        item.Resource + ":" + item.Action,
			Description: item.Description,
		})
	}
	return out
}

func fallback(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func sessionStatus(item identityread.SessionReadModel, now time.Time) string {
	if item.RevokedAt != nil {
		return "revogada"
	}
	if !now.Before(item.ExpiresAt) {
		return "expirada"
	}
	return "ativa"
}
