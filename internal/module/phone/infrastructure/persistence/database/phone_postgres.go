// Package database
package database

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/module/phone/domain"
)

type Phone struct {
	pool postgres.PgxPool
}

func NewPhone(pool postgres.PgxPool) *Phone {
	return &Phone{
		pool,
	}
}

func (p *Phone) Create(ctx context.Context, phone *domain.Phone) error {
	fmt.Println(phone.CreatedAt(), phone.UpdatedAt())
	query := `INSERT INTO phones (number, kind, department, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING idphones`
	var id uint
	if err := p.pool.QueryRow(ctx, query,
		phone.Number(),
		phone.Kind(),
		phone.Department(),
		phone.CreatedAt(),
		phone.UpdatedAt()).Scan(&id); err != nil {
		return err
	}
	return phone.SetID(id)
}
