// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/readmodel"
)

type CompanyListFilter struct {
	Name string
	CNPJ string
}

type CompanyReadRepository interface {
	List(ctx context.Context, filter CompanyListFilter, after, before *string, limit int, baseURL string) (*readmodel.CompanyListReadModel, error)
}
