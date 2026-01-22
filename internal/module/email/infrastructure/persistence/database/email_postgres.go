// Package database
package database

import (
	"context"

	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/module/email/domain"
)

type Email struct {
	pool postgres.PgxPool
}

func NewEmail(pool postgres.PgxPool) *Email {
	return &Email{
		pool,
	}
}

func (e *Email) Create(ctx context.Context, email *domain.Email) error {
	query := `INSERT INTO emails (address, created_at, updated_at)
		VALUES ($1, $2, $3) RETURNING idemails`
	var id uint
	if err := e.pool.QueryRow(ctx, query,
		email.Address(),
		email.CreatedAt(),
		email.UpdatedAt()).Scan(&id); err != nil {
		return err
	}
	return email.SetID(id)
}
