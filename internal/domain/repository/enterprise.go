// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type IEnterprise interface {
	Create(ctx context.Context, enterprise entity.Enterprise) (id string, err error)
}
