package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestAddressRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	address, err := addrDomain.NewAddressBuilder().
		WithZIP("79000-000").
		WithTitle("Matriz").
		WithStreet("Rua A").
		WithNumber(42).
		WithComplement("Sala 1").
		WithReference("Esquina").
		WithNeighborhood("Centro").
		WithCity("Campo Grande").
		WithState("MS").
		WithCountry("BR").
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
	require.NoError(t, err)

	mock.ExpectQuery(`INSERT INTO addresses`).
		WithArgs(
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
		).
		WillReturnRows(pgxmock.NewRows([]string{"idaddresses"}).AddRow(int64(10)))

	repo := NewAddressRepository(mock)
	err = repo.Create(ctx, address)
	require.NoError(t, err)
	require.Equal(t, int64(10), address.ID().Value())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddressRepositoryFindByIDCompany(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`SELECT\s+a.idaddresses`).
		WithArgs(int64(55)).
		WillReturnRows(pgxmock.NewRows([]string{
			"idaddresses", "zip", "title", "street", "number", "complement", "reference",
			"neighborhood", "city", "state", "country", "created_at", "updated_at",
		}).AddRow(
			int64(1), "79000-000", "Matriz", "Rua A", uint(42), "Sala 1", "Esquina",
			"Centro", "Campo Grande", "MS", "BR", now, now,
		))

	repo := NewAddressRepository(mock)
	got, err := repo.FindByIDCompany(ctx, 55)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, int64(1), got.ID().Value())
	require.Equal(t, "79000-000", got.ZIP().Value())
	require.Equal(t, "Campo Grande", got.City().Value())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddressRepositoryFindByIDCompanyNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectQuery(`SELECT\s+a.idaddresses`).
		WithArgs(int64(99)).
		WillReturnError(pgx.ErrNoRows)

	repo := NewAddressRepository(mock)
	got, err := repo.FindByIDCompany(ctx, 99)
	require.NoError(t, err)
	require.Nil(t, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddressRepositoryFindByIDCompanyInvalidID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectQuery(`SELECT\s+a.idaddresses`).
		WithArgs(int64(0)).
		WillReturnError(&pgconn.PgError{Code: "22P02"})

	repo := NewAddressRepository(mock)
	got, err := repo.FindByIDCompany(ctx, 0)
	require.Nil(t, got)
	require.Error(t, err)
	fieldErr, ok := err.(*sharedDomain.FieldError)
	require.True(t, ok)
	require.Equal(t, "id_company", fieldErr.Field)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddressRepositoryUpdate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	address, err := addrDomain.NewAddressBuilder().
		WithID(10).
		WithZIP("79000-000").
		WithTitle("Matriz").
		WithStreet("Rua A").
		WithNumber(42).
		WithComplement("Sala 1").
		WithReference("Esquina").
		WithNeighborhood("Centro").
		WithCity("Campo Grande").
		WithState("MS").
		WithCountry("BR").
		WithUpdatedAt(now).
		Build()
	require.NoError(t, err)

	mock.ExpectExec(`UPDATE addresses`).
		WithArgs(
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
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	repo := NewAddressRepository(mock)
	err = repo.Update(ctx, *address)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
