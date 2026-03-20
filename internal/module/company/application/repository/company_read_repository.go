// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/readmodel"
)

type CompanyReadRepository interface {
	List(ctx context.Context, after, before *string, limit int, baseURL string) (*readmodel.CompanyListReadModel, error)
}
