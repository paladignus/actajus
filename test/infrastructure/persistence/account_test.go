// Package persistence
package persistence

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/infrastructure/database"
	"github.com/paladignus/actajus/internal/infrastructure/persistence"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount_FindUserAccountByCPF(t *testing.T) {
	ctx := context.Background()
	t.Run("should return user account with roles", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()
		db := &database.DB{Pool: mock}
		repo := persistence.NewAccount(db)
		cpf := "12345678900"
		expectedIDUser := "user-123"
		expectedIDAccount := "account-456"
		expectedFirstName := "John"
		expectedLastName := "Doe"
		expectedEmail := "john.doe@example.com"
		expectedRoles := []string{"admin", "user"}
		rows := pgxmock.NewRows([]string{"idpeople", "first_name", "last_name", "address", "idaccounts"}).
			AddRow(expectedIDUser, expectedFirstName, expectedLastName, expectedEmail, expectedIDAccount)
		mock.ExpectQuery(`SELECT p.idpeople, p.first_name, p.last_name, e.address, a.idaccounts FROM people p`).
			WithArgs(cpf).
			WillReturnRows(rows)
		rolesRows := pgxmock.NewRows([]string{"name"}).
			AddRow("admin").
			AddRow("user")
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(expectedIDAccount).
			WillReturnRows(rolesRows)
		result, err := repo.FindUserAccountByCPF(ctx, cpf)
		assert.NoError(t, err)
		assert.Equal(t, expectedIDUser, result.IDUser)
		assert.Equal(t, expectedFirstName, result.FirstName)
		assert.Equal(t, expectedLastName, result.LastName)
		assert.Equal(t, expectedEmail, result.Email)
		assert.Equal(t, expectedRoles, result.Roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return user account not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()
		db := &database.DB{Pool: mock}
		repo := persistence.NewAccount(db)
		cpf := "99999999999"
		mock.ExpectQuery(`SELECT p.idpeople, p.first_name, p.last_name, e.address, a.idaccounts FROM people p`).
			WithArgs(cpf).
			WillReturnError(pgx.ErrNoRows)
		result, err := repo.FindUserAccountByCPF(ctx, cpf)
		assert.Error(t, err)
		assert.Equal(t, exception.ErrUserNotFound, err)
		assert.Equal(t, dto.AuthenticateOutput{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error from database", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()
		db := &database.DB{Pool: mock}
		repo := persistence.NewAccount(db)
		cpf := "12345678900"
		expectedErr := errors.New("database connection error")
		mock.ExpectQuery(`SELECT p.idpeople, p.first_name, p.last_name, e.address, a.idaccounts FROM people p`).
			WithArgs(cpf).
			WillReturnError(expectedErr)
		result, err := repo.FindUserAccountByCPF(ctx, cpf)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, dto.AuthenticateOutput{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should error from roles query", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()
		db := &database.DB{Pool: mock}
		repo := persistence.NewAccount(db)
		cpf := "12345678900"
		expectedErr := errors.New("roles query error")
		rows := pgxmock.NewRows([]string{"idpeople", "first_name", "last_name", "address", "idaccounts"}).
			AddRow("user-123", "John", "Doe", "john@example.com", "account-456")
		mock.ExpectQuery(`SELECT p.idpeople, p.first_name, p.last_name, e.address, a.idaccounts FROM people p`).
			WithArgs(cpf).
			WillReturnRows(rows)
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs("account-456").
			WillReturnError(expectedErr)
		result, err := repo.FindUserAccountByCPF(ctx, cpf)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, dto.AuthenticateOutput{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
