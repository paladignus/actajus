// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
)

type CompanyReadRepository interface {
	List(ctx context.Context, after, before *string, limit int, baseURL string) (*dto.CompanyListReadModel, error)
}
