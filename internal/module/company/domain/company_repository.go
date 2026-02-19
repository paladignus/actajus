// Package domain
package domain

import (
	"context"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, company *Company) error
	FindByCNPJ(ctx context.Context, cnpj string) (*Company, error)
	FindByID(ctx context.Context, id uint) (*Company, error)
}
