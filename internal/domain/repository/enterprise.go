// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type IEnterprise interface {
	Create(ctx context.Context, enterprise entity.Enterprise) (id uint, err error)
	Update(ctx context.Context, enterprise entity.Enterprise) error
	GetAll(ctx context.Context) ([]entity.Enterprise, error)
}
