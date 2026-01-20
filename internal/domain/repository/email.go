// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type IEmail interface {
	Create(ctx context.Context, email entity.Email) (id uint, err error)
	Update(ctx context.Context, email entity.Email) error
}
