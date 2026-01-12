// Package unitofwork
package unitofwork

import (
	"context"

	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type IUnitOfWork interface {
	Begin(context.Context) error
	Commit(context.Context) error
	Rollback(context.Context) error
	GetPgxPool() database.PgxPool
}
