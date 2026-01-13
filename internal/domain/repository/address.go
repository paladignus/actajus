// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type IAddress interface {
	Create(ctx context.Context, address entity.Address) (string, error)
	// Get(ctx context.Context, idAddress uint) (entity.Address, error)
	// Update(ctx context.Context, address entity.Address) (entity.Address, error)
	// Delete(ctx context.Context, idAddress uint) error
}
