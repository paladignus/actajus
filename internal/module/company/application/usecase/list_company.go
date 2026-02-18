// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type ListCompanies struct {
	list domain.CompanyRepository
}

func NewListCompanies(list domain.CompanyRepository) ListCompanies {
	return ListCompanies{list}
}

func (l ListCompanies) Execute(ctx context.Context, after, before *string, limit int, baseURL string) (*dto.CompanyListReadModel, error) {
	return l.list.List(ctx, after, before, limit, baseURL)
}
