// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type IPhone interface {
	Create(ctx context.Context, phone entity.Phone) (id uint, err error)
	Update(ctx context.Context, phone entity.Phone) error
}
