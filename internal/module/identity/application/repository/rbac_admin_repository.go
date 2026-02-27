// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/domain"
)

type RoleUserAdminRepository interface {
	AssignRole(ctx context.Context, idUser domain.IDUser, idRole int16, assignedBy int64) error
	RemoveRole(ctx context.Context, idUser domain.IDUser, idRole int16) error
	ListUserIDsByRole(ctx context.Context, idRole int16) ([]domain.IDUser, error)
}

type PermissionRoleAdminRepository interface {
	GrantPermission(ctx context.Context, idRole int16, idPermission int16) error
	RevokePermission(ctx context.Context, idRole int16, idPermission int16) error
}
