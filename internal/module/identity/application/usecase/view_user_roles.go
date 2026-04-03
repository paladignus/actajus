// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type ViewUserRoles struct {
	query repository.CatalogQueryRepository
}

func NewViewUserRoles(query repository.CatalogQueryRepository) ViewUserRoles {
	return ViewUserRoles{query: query}
}

func (uc ViewUserRoles) Execute(ctx context.Context, idUser int64) ([]readmodel.UserRoleAssignmentReadModel, error) {
	return uc.query.ListRolesByUser(ctx, idUser)
}
