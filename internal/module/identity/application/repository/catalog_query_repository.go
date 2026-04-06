// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
)

type CatalogUserFilter struct {
	Query string
	Limit int
}

type CatalogPermissionFilter struct {
	Query string
	Limit int
}

type CatalogQueryRepository interface {
	ListUsers(ctx context.Context, filter CatalogUserFilter) ([]readmodel.UserListItemReadModel, error)
	ListPermissions(ctx context.Context, filter CatalogPermissionFilter) ([]readmodel.PermissionListItemReadModel, error)
	ListRoles(ctx context.Context) ([]readmodel.RoleListItemReadModel, error)
	ListRolesByUser(ctx context.Context, idUser int64) ([]readmodel.UserRoleAssignmentReadModel, error)
	ListPermissionsByRole(ctx context.Context, idRole int16) ([]readmodel.RolePermissionAssignmentReadModel, error)
}
