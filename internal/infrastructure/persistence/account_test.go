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
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount_GetEmailByCPF(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewAccount(db)
	t.Run("should return email by cpf", func(t *testing.T) {
		cpf := "11144477735"
		expectedEmail := "email@example.com.br"
		rows := pgxmock.NewRows([]string{"address"}).AddRow(expectedEmail)
		mock.ExpectQuery(`SELECT e.address FROM emails e`).
			WithArgs(cpf).
			WillReturnRows(rows)
		email, err := repo.GetEmailByCPF(ctx, cpf)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmail, email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when email not found", func(t *testing.T) {
		cpf := "00000000000"
		mock.ExpectQuery(`SELECT e.address FROM emails e`).
			WithArgs(cpf).
			WillReturnError(pgx.ErrNoRows)
		email, err := repo.GetEmailByCPF(ctx, cpf)
		assert.Error(t, err)
		assert.Equal(t, exception.ErrEmailNotFound, err)
		assert.Equal(t, "", email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAccount_FindUserAccountByCPF(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewAccount(db)
	t.Run("should return user account with roles", func(t *testing.T) {
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

func TestAccount_ValidatePassword(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewAccount(db)
	t.Run("should validate password", func(t *testing.T) {
		idPeople := "user-123"
		password := "correct-password"
		rows := pgxmock.NewRows([]string{"valid"}).AddRow(true)
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs(idPeople, password).
			WillReturnRows(rows)
		err = repo.ValidatePassword(ctx, idPeople, password)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should error of invalid credentials", func(t *testing.T) {
		idPeople := "user-123"
		password := "wrong-password"
		rows := pgxmock.NewRows([]string{"valid"}).AddRow(false)
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs(idPeople, password).
			WillReturnRows(rows)
		err = repo.ValidatePassword(ctx, idPeople, password)
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidCredentials, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error from database", func(t *testing.T) {
		idPeople := "user-123"
		password := "any-password"
		expectedErr := errors.New("database error")
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs(idPeople, password).
			WillReturnError(expectedErr)
		err = repo.ValidatePassword(ctx, idPeople, password)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAccount_GetRolesByAccountID(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewAccount(db)
	t.Run("should return roles", func(t *testing.T) {
		idAccount := "account-123"
		expectedRoles := []string{"admin", "manager", "user"}
		rows := pgxmock.NewRows([]string{"name"}).
			AddRow("admin").
			AddRow("manager").
			AddRow("user")
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(idAccount).
			WillReturnRows(rows)
		roles, err := repo.GetRolesByAccountID(ctx, idAccount)
		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should empty roles", func(t *testing.T) {
		idAccount := "account-123"
		rows := pgxmock.NewRows([]string{"name"})
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(idAccount).
			WillReturnRows(rows)
		roles, err := repo.GetRolesByAccountID(ctx, idAccount)
		assert.NoError(t, err)
		assert.Nil(t, roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error", func(t *testing.T) {
		idAccount := "account-123"
		expectedErr := errors.New("query error")
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(idAccount).
			WillReturnError(expectedErr)
		roles, err := repo.GetRolesByAccountID(ctx, idAccount)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should error from database", func(t *testing.T) {
		idAccount := "account-123"
		rows := pgxmock.NewRows([]string{"name"}).
			AddRow("admin").
			AddRow(nil).
			RowError(1, errors.New("scan error"))
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(idAccount).
			WillReturnRows(rows)
		roles, err := repo.GetRolesByAccountID(ctx, idAccount)
		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAccount_GetPermissionsByRoleID(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewAccount(db)
	t.Run("should return permissions", func(t *testing.T) {
		roleID := 1
		expectedPermissions := []dto.Permission{
			{Resource: "users", Action: "create"},
			{Resource: "users", Action: "read"},
			{Resource: "posts", Action: "delete"},
		}
		rows := pgxmock.NewRows([]string{"resource", "action"}).
			AddRow("users", "create").
			AddRow("users", "read").
			AddRow("posts", "delete")
		mock.ExpectQuery(`SELECT p.resource, p.action FROM permissions p`).
			WithArgs(roleID).
			WillReturnRows(rows)
		permissions, err := repo.GetPermissionsByRoleID(ctx, roleID)
		assert.NoError(t, err)
		assert.Equal(t, expectedPermissions, permissions)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return empty permissions", func(t *testing.T) {
		roleID := 1
		rows := pgxmock.NewRows([]string{"resource", "action"})
		mock.ExpectQuery(`SELECT p.resource, p.action FROM permissions p`).
			WithArgs(roleID).
			WillReturnRows(rows)
		permissions, err := repo.GetPermissionsByRoleID(ctx, roleID)
		assert.NoError(t, err)
		assert.Nil(t, permissions)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error to query", func(t *testing.T) {
		roleID := 1
		expectedErr := errors.New("query error")
		mock.ExpectQuery(`SELECT p.resource, p.action FROM permissions p`).
			WithArgs(roleID).
			WillReturnError(expectedErr)
		permissions, err := repo.GetPermissionsByRoleID(ctx, roleID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, permissions)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error to scan", func(t *testing.T) {
		roleID := 1
		rows := pgxmock.NewRows([]string{"resource", "action"}).
			AddRow("users", "create").
			AddRow(nil, nil).
			RowError(1, errors.New("scan error"))
		mock.ExpectQuery(`SELECT p.resource, p.action FROM permissions p`).
			WithArgs(roleID).
			WillReturnRows(rows)
		permissions, err := repo.GetPermissionsByRoleID(ctx, roleID)
		assert.Error(t, err)
		assert.Nil(t, permissions)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
