// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
	"github.com/paladignus/actajus/internal/module/company/infrastructure/persistence/database"
)

type ListCompanies struct {
	repository domain.CompanyReadRepository
}

func NewListCompanies(repository database.CompanyReadRepository) ListCompanies {
	return ListCompanies{repository}
}

func (l ListCompanies) Execute(ctx context.Context, after, before *string, limit int, baseURL string) (*dto.CompanyListReadModel, error) {
	return l.repository.List(ctx, after, before, limit, baseURL)
}
