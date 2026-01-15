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

func (p Phone) Update(ctx context.Context, phone entity.Phone) error {
	sql := `UPDATE phones SET number = $1, kind = $2, department = $3, updated_at = now() WHERE idphones = $4;`
	_, err := p.tx.Exec(ctx, sql, phone.Number, phone.Kind, phone.Department, phone.IDPhone)
	if err != nil {
		return fmt.Errorf("database error while updating phone for user ID %d: %w", phone.IDPhone, err)
	}
	return nil
}
