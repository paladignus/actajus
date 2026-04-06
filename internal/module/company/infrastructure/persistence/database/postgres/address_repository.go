// Package postgres
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	sharedpostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type AddressRepository struct {
	exec sharedpostgres.Executor
}

func NewAddressRepository(exec sharedpostgres.Executor) AddressRepository {
	return AddressRepository{exec: exec}
}

func (r AddressRepository) Create(ctx context.Context, address *addrDomain.Address) error {
	query := `INSERT INTO addresses
		(zip, title, street, number, complement, reference, neighborhood, city, state, country, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING idaddresses`
	var id int64
	if err := r.exec.QueryRow(ctx, query,
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

func (r AddressRepository) Update(ctx context.Context, address addrDomain.Address) error {
	query := `
		UPDATE addresses
		SET zip = $1, title = $2, street = $3, number = $4, complement = $5,
		reference = $6, neighborhood = $7, city = $8, state = $9, country = $10, updated_at = $11
		WHERE idaddresses = $12`
	_, err := r.exec.Exec(ctx, query,
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
		address.UpdatedAt(),
		address.ID().Value(),
	)
	return err
}

func (r AddressRepository) FindByIDCompany(ctx context.Context, idCompany int64) (*addrDomain.Address, error) {
	query := `
		SELECT
			a.idaddresses, a.zip, a.title, a.street, a.number, a.complement, a.reference,
			a.neighborhood, a.city, a.state, a.country, a.created_at, a.updated_at
		FROM addresses a
		LEFT JOIN company_address ca ON a.idaddresses = ca.id_addresses AND ca.ended_at IS NULL
		WHERE id_companies = $1`
	var (
		id           int64
		zip          string
		title        string
		street       string
		number       uint
		complement   string
		reference    string
		neighborhood string
		city         string
		state        string
		country      string
		createdAt    time.Time
		updatedAt    time.Time
	)
	err := r.exec.QueryRow(ctx, query, idCompany).
		Scan(&id, &zip, &title, &street, &number, &complement, &reference, &neighborhood, &city, &state, &country, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, sharedDomain.NewFieldError("id_company", "idCompany")
		}
		return nil, err
	}
	return addrDomain.NewAddressBuilder().
		WithID(id).
		WithZIP(zip).
		WithTitle(title).
		WithStreet(street).
		WithNumber(number).
		WithComplement(complement).
		WithReference(reference).
		WithNeighborhood(neighborhood).
		WithCity(city).
		WithState(state).
		WithCountry(country).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()
}
