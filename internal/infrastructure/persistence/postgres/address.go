// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/domain/entity"
)

type Address struct {
	tx pgx.Tx
}

func NewAddress(tx pgx.Tx) Address {
	return Address{tx}
}

func (a Address) Create(ctx context.Context, input entity.Address) (id string, err error) {
	sql := `INSERT INTO addresses 
		(zip, title, street, number, complement, neighborhood, city, state, country)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning idaddresses;`
	if err := a.tx.QueryRow(ctx, sql,
		input.Zip, input.Title, input.Street, input.Number, input.Complement,
		input.Neighborhood, input.City, input.State, input.Country,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("database error while saving address: %w", err)
	}
	return id, nil
}
