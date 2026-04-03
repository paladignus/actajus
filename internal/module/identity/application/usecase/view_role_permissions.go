// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type ViewRolePermissions struct {
	query repository.CatalogQueryRepository
}

func NewViewRolePermissions(query repository.CatalogQueryRepository) ViewRolePermissions {
	return ViewRolePermissions{query: query}
}

func (uc ViewRolePermissions) Execute(ctx context.Context, idRole int16) ([]readmodel.RolePermissionAssignmentReadModel, error) {
	return uc.query.ListPermissionsByRole(ctx, idRole)
}
