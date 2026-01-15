// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type Address struct {
	tx PgxPool
}

func NewAddress(tx PgxPool) Address {
	return Address{tx}
}

func (a Address) Create(ctx context.Context, input entity.Address) (id uint, err error) {
	sql := `INSERT INTO addresses 
		(zip, title, street, number, complement, reference, neighborhood, city, state, country)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning idaddresses;`
	if err := a.tx.QueryRow(ctx, sql,
		input.Zip, input.Title, input.Street, input.Number, input.Complement, input.Reference,
		input.Neighborhood, input.City, input.State, input.Country,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("database error while saving address: %w", err)
	}
	return id, nil
}

func (a Address) Update(ctx context.Context, input entity.Address) error {
	sql := `UPDATE addresses
		SET zip = $1, title = $2, street = $3, number = $4, complement = $5, reference = $6,
	neighborhood = $7, city = $8, state = $9, country = $10, updated_at = now() WHERE idaddresses = $11;`
	if _, err := a.tx.Exec(ctx, sql,
		input.Zip, input.Title, input.Street, input.Number, input.Complement, input.Reference,
		input.Neighborhood, input.City, input.State, input.Country, input.IDAddress,
	); err != nil {
		return fmt.Errorf("database error while update address: %w", err)
	}
	return nil
}
