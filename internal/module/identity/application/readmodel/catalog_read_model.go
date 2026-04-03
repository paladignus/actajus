// Package readmodel
package readmodel

import "time"

type UserListItemReadModel struct {
	ID          int64      `json:"id"`
	FullName    string     `json:"full_name"`
	Email       string     `json:"email"`
	IsBlocked   bool       `json:"is_blocked"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type PermissionListItemReadModel struct {
	ID          int16     `json:"id"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type RoleListItemReadModel struct {
	ID          int16     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserRoleAssignmentReadModel struct {
	IDUser     int64  `json:"id_user"`
	IDRole     int16  `json:"id_role"`
	RoleName   string `json:"role_name"`
	AssignedBy int64  `json:"assigned_by"`
}

type RolePermissionAssignmentReadModel struct {
	IDRole       int16  `json:"id_role"`
	IDPermission int16  `json:"id_permission"`
	Resource     string `json:"resource"`
	Action       string `json:"action"`
	Description  string `json:"description"`
}
