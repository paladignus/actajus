// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type Phone struct {
	tx PgxPool
}

func NewPhone(tx PgxPool) Phone {
	return Phone{tx}
}

func (p Phone) Create(ctx context.Context, phone entity.Phone) (id uint, err error) {
	sql := `INSERT INTO phones (number, kind, department) VALUES ($1, $2, $3) RETURNING idphones;`
	if err := p.tx.QueryRow(ctx, sql, phone.Number, phone.Kind, phone.Department).Scan(&id); err != nil {
		return 0, fmt.Errorf("database error while saving phone for user ID %d: %w", id, err)
	}
	return id, nil
}
