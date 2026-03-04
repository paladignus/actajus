// Package repository
package repository

import (
	"context"
)

type RoleUserAdminRepository interface {
	AssignRole(ctx context.Context, uid int64, rid int16, assignedBy int64) error
	RemoveRole(ctx context.Context, uid int64, rid int16) error
	ListUserIDsByRole(ctx context.Context, idRole int16) ([]int64, error)
}

type PermissionRoleAdminRepository interface {
	GrantPermission(ctx context.Context, rid int16, pid int16) error
	RevokePermission(ctx context.Context, rid int16, pid int16) error
}
