// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type ListRoles struct {
	query repository.CatalogQueryRepository
}

func NewListRoles(query repository.CatalogQueryRepository) ListRoles {
	return ListRoles{query: query}
}

func (uc ListRoles) Execute(ctx context.Context) ([]readmodel.RoleListItemReadModel, error) {
	return uc.query.ListRoles(ctx)
}
