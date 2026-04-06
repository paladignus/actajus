// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type ListPermissions struct {
	query repository.CatalogQueryRepository
}

func NewListPermissions(query repository.CatalogQueryRepository) ListPermissions {
	return ListPermissions{query: query}
}

func (uc ListPermissions) Execute(ctx context.Context, filter repository.CatalogPermissionFilter) ([]readmodel.PermissionListItemReadModel, error) {
	return uc.query.ListPermissions(ctx, filter)
}
