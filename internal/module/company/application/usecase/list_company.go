// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/readmodel"
	"github.com/paladignus/actajus/internal/module/company/application/repository"
)

type ListCompanies struct {
	repository repository.CompanyReadRepository
}

func NewListCompanies(repository repository.CompanyReadRepository) ListCompanies {
	return ListCompanies{repository}
}

func (l ListCompanies) Execute(ctx context.Context, filter repository.CompanyListFilter, after, before *string, limit int, baseURL string) (*readmodel.CompanyListReadModel, error) {
	return l.repository.List(ctx, filter, after, before, limit, baseURL)
}
