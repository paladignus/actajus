// Package domain
package domain

import "context"

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	FindByCNPJ(ctx context.Context, cnpj string) (*Company, error)
}
