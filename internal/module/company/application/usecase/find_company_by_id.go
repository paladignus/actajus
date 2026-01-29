// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/domain"
)

type FindByID struct {
	repository domain.CompanyRepository
}

func NewFindByID(repository domain.CompanyRepository) FindByID {
	return FindByID{repository}
}

func (f FindByID) Execute(ctx context.Context, id uint) (*domain.Company, error) {
	return f.repository.FindByID(ctx, id)
}
