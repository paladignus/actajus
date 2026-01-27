// Package domain
package domain

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	FindByCNPJ(ctx context.Context, cnpj string) (*Company, error)
	List(ctx context.Context, page, pageSize uint) (*dto.CompanyListReadModel, error)
}
