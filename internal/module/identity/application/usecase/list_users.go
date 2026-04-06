// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type ListUsers struct {
	query repository.CatalogQueryRepository
}

func NewListUsers(query repository.CatalogQueryRepository) ListUsers {
	return ListUsers{query: query}
}

func (uc ListUsers) Execute(ctx context.Context, filter repository.CatalogUserFilter) ([]readmodel.UserListItemReadModel, error) {
	return uc.query.ListUsers(ctx, filter)
}
