// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/domain"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *domain.Company) error
	Update(ctx context.Context, company *domain.Company) error
	Delete(ctx context.Context, company *domain.Company) error
	FindByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error)
	FindByID(ctx context.Context, id int64) (*domain.Company, error)
}
