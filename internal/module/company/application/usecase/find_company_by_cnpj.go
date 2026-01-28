// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type FindByCNPJ struct {
	repository domain.CompanyRepository
}

func NewFindByCNPJ(repository domain.CompanyRepository) FindByCNPJ {
	return FindByCNPJ{repository: repository}
}

func (f FindByCNPJ) Execute(ctx context.Context, cnpj string) (*dto.CompanyReadModel, error) {
	return f.repository.FindByCNPJ(ctx, cnpj)
}
