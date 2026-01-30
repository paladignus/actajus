// Package domain
package domain

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, company *Company) error
	FindByCNPJ(ctx context.Context, cnpj string) (*Company, error)
	FindByID(ctx context.Context, id uint) (*Company, error)
	List(ctx context.Context, page, pageSize uint) (*dto.CompanyListReadModel, error)
}
