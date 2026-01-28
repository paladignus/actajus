// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type FindByID struct {
	repository domain.CompanyRepository
}

func NewFindByID(repository domain.CompanyRepository) FindByID {
	return FindByID{repository: repository}
}

func (f FindByID) Execute(ctx context.Context, cnpj string) (*dto.CompanyReadModel, error) {
	return f.repository.FindByID(ctx, cnpj)
}
