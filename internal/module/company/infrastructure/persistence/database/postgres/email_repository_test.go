package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestEmailRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	email, err := emailDomain.NewEmailBuilder().
		WithAddress("contato@acme.com").
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
	require.NoError(t, err)

	mock.ExpectQuery(`INSERT INTO emails`).
		WithArgs(email.Address(), email.CreatedAt(), email.UpdatedAt()).
		WillReturnRows(pgxmock.NewRows([]string{"idemails"}).AddRow(int64(21)))

	repo := NewEmailRepository(mock)
	err = repo.Create(ctx, email)
	require.NoError(t, err)
	require.Equal(t, int64(21), email.ID().Value())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestEmailRepositoryFindByIDCompany(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`SELECT\s+e.idemails`).
		WithArgs(int64(55)).
		WillReturnRows(pgxmock.NewRows([]string{
			"idemails", "address", "created_at", "updated_at",
		}).AddRow(int64(1), "contato@acme.com", now, now))

	repo := NewEmailRepository(mock)
	got, err := repo.FindByIDCompany(ctx, 55)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, int64(1), got.ID().Value())
	require.Equal(t, "contato@acme.com", got.Address().Value())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestEmailRepositoryFindByIDCompanyErrorPaths(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		mock.ExpectQuery(`SELECT\s+e.idemails`).
			WithArgs(int64(99)).
			WillReturnError(pgx.ErrNoRows)

		repo := NewEmailRepository(mock)
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

		mock.ExpectQuery(`SELECT\s+e.idemails`).
			WithArgs(int64(0)).
			WillReturnError(&pgconn.PgError{Code: "22P02"})

		repo := NewEmailRepository(mock)
		got, err := repo.FindByIDCompany(ctx, 0)
		require.Nil(t, got)
		require.Error(t, err)
		fieldErr, ok := err.(*sharedDomain.FieldError)
		require.True(t, ok)
		require.Equal(t, "id_company", fieldErr.Field)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmailRepositoryUpdateAndDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	deletedAt := now.Add(time.Minute)
	email, err := emailDomain.NewEmailBuilder().
		WithID(21).
		WithAddress("contato@acme.com").
		WithUpdatedAt(now).
		WithDeletedAt(&deletedAt).
		Build()
	require.NoError(t, err)

	mock.ExpectExec(`UPDATE emails SET address =`).
		WithArgs(email.Address(), email.UpdatedAt(), email.ID().Value()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec(`UPDATE emails SET updated_at =`).
		WithArgs(email.UpdatedAt(), email.DeletedAt(), email.ID().Value()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	repo := NewEmailRepository(mock)
	err = repo.Update(ctx, *email)
	require.NoError(t, err)
	err = repo.Delete(ctx, *email)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
