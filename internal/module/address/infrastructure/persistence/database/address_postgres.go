// Package database
package database

import (
	"context"

	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/module/address/domain"
)

type Address struct {
	pool postgres.PgxPool
}

func NewAddress(pool postgres.PgxPool) *Address {
	return &Address{
		pool,
	}
}

func (a Address) FindByCEP(ctx context.Context, cep string) (*domain.Address, error) {
	return nil, nil
}

func (a Address) Create(ctx context.Context, address *domain.Address) error {
	query := `INSERT INTO addresses
		(zip, title, street, number, complement, reference, neighborhood, city, state, country, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING idaddresses`
	var id uint
	if err := a.pool.QueryRow(ctx, query,
		address.ZIP().Value(),
		address.Title().Value(),
		address.Street().Value(),
		address.Number(),
		address.Complement().Value(),
		address.Reference().Value(),
		address.Neighborhood().Value(),
		address.City().Value(),
		address.State().Value(),
		address.Country().Value(),
		address.CreatedAt(),
		address.UpdatedAt(),
	).Scan(&id); err != nil {
		return err
	}
	return address.SetID(id)
}
