package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestPhoneRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	phone, err := phoneDomain.NewPhoneBuilder().
		WithNumber("+5567999999999").
		WithKind("commercial").
		WithDepartment("sales").
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
	require.NoError(t, err)

	mock.ExpectQuery(`INSERT INTO phones`).
		WithArgs(
			phone.Number(),
			phone.Kind(),
			phone.Department(),
			phone.CreatedAt(),
			phone.UpdatedAt(),
		).
		WillReturnRows(pgxmock.NewRows([]string{"idphones"}).AddRow(int64(12)))

	repo := NewPhoneRepository(mock)
	err = repo.Create(ctx, phone)
	require.NoError(t, err)
	require.Equal(t, int64(12), phone.ID().Value())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPhoneRepositoryFindByIDCompany(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`SELECT\s+p.idphones`).
		WithArgs(int64(55)).
		WillReturnRows(pgxmock.NewRows([]string{
			"idphones", "number", "kind", "department", "created_at", "updated_at",
		}).AddRow(int64(1), "+5567999999999", "commercial", "sales", now, now))

	repo := NewPhoneRepository(mock)
	got, err := repo.FindByIDCompany(ctx, 55)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, int64(1), got.ID().Value())
	require.Equal(t, "+5567999999999", got.Number().Value())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPhoneRepositoryFindByIDCompanyErrorPaths(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		mock.ExpectQuery(`SELECT\s+p.idphones`).
			WithArgs(int64(99)).
			WillReturnError(pgx.ErrNoRows)

		repo := NewPhoneRepository(mock)
		got, err := repo.FindByIDCompany(ctx, 99)
		require.NoError(t, err)
		require.Nil(t, got)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid id", func(t *testing.T) {
		ctx := context.Background()
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		mock.ExpectQuery(`SELECT\s+p.idphones`).
			WithArgs(int64(0)).
			WillReturnError(&pgconn.PgError{Code: "22P02"})

		repo := NewPhoneRepository(mock)
		got, err := repo.FindByIDCompany(ctx, 0)
		require.Nil(t, got)
		require.Error(t, err)
		fieldErr, ok := err.(*sharedDomain.FieldError)
		require.True(t, ok)
		require.Equal(t, "id_company", fieldErr.Field)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPhoneRepositoryUpdateAndDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	deletedAt := now.Add(time.Minute)
	phone, err := phoneDomain.NewPhoneBuilder().
		WithID(12).
		WithNumber("+5567999999999").
		WithKind("commercial").
		WithDepartment("sales").
		WithUpdatedAt(now).
		WithDeletedAt(&deletedAt).
		Build()
	require.NoError(t, err)

	mock.ExpectExec(`UPDATE phones SET number =`).
		WithArgs(phone.Number(), phone.Kind(), phone.Department(), phone.UpdatedAt(), phone.ID().Value()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec(`UPDATE phones SET updated_at =`).
		WithArgs(phone.UpdatedAt(), phone.DeletedAt(), phone.ID().Value()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	repo := NewPhoneRepository(mock)
	err = repo.Update(ctx, *phone)
	require.NoError(t, err)
	err = repo.Delete(ctx, *phone)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
