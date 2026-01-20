// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type Email struct {
	db PgxPool
}

func NewEmail(db PgxPool) Email {
	return Email{db}
}

func (e Email) Create(ctx context.Context, email entity.Email) (id uint, err error) {
	sql := `INSERT INTO emails (address) VALUES ($1) RETURNING idemails;`
	if err = e.db.QueryRow(ctx, sql, email.Address.Value()).Scan(&id); err != nil {
		return 0, fmt.Errorf("database error while saving email: %w", err)
	}
	return id, nil
}

func (e Email) Update(ctx context.Context, email entity.Email) error {
	sql := `UPDATE emails SET address = $1, updated_at = now() WHERE idemails = $2;`
	if _, err := e.db.Exec(ctx, sql, email.Address.Value(), email.IDEmails); err != nil {
		return fmt.Errorf("database error while update email: %w", err)
	}
	return nil
}
