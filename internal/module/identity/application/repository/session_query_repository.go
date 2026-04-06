// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
)

type SessionQueryFilter struct {
	UserEmail string
	Limit     int
}

type SessionQueryRepository interface {
	GetByID(ctx context.Context, id int64) (*readmodel.SessionReadModel, error)
	ListByUser(ctx context.Context, uid int64) ([]readmodel.SessionReadModel, error)
	ListAll(ctx context.Context, filter SessionQueryFilter) ([]readmodel.SessionReadModel, error)
}
