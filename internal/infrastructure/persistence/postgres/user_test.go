// Package postgres
package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_AuthenticationByCPF(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	// db := &database.DB{Pool: mock}
	repo := NewUser(mock)
	input := command.SignInCommand{
		CPF:      "12345678900",
		Password: "Password123@",
	}
	expectedIDUser := 123
	expectedFirstName := "John"
	expectedLastName := "Doe"
	expectedEmail := "john.doe@example.com"
	t.Run("should return user account with roles", func(t *testing.T) {
		expectedRoles := []string{"admin", "user"}
		rows := pgxmock.NewRows([]string{"idusers", "first_name", "last_name", "address"}).
			AddRow(expectedIDUser, expectedFirstName, expectedLastName, expectedEmail)
		mock.ExpectQuery(`SELECT u.idusers, p.first_name, p.last_name, e.address FROM documents d`).
			WithArgs(input.CPF, input.Password).
			WillReturnRows(rows)
		rolesRows := pgxmock.NewRows([]string{"name"}).
			AddRow("admin").
			AddRow("user")
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(expectedIDUser).
			WillReturnRows(rolesRows)
		result, err := repo.AuthenticationByCPF(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, expectedIDUser, result.IDUser)
		assert.Equal(t, expectedFirstName, result.FirstName)
		assert.Equal(t, expectedLastName, result.LastName)
		assert.Equal(t, expectedEmail, result.Email)
		assert.Equal(t, expectedRoles, result.Roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return user account not found", func(t *testing.T) {
		input.CPF = "99999999999"
		mock.ExpectQuery(`SELECT u.idusers, p.first_name, p.last_name, e.address FROM documents d`).
			WithArgs(input.CPF, input.Password).
			WillReturnError(pgx.ErrNoRows)
		result, err := repo.AuthenticationByCPF(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "authentication failed for CPF 99999999999")
		assert.ErrorContains(t, err, "invalid credentials")
		assert.Equal(t, readmodel.SignInReadModel{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error from database", func(t *testing.T) {
		input.CPF = "12345678900"
		expectedErr := errors.New("database connection error")
		mock.ExpectQuery(`SELECT u.idusers, p.first_name, p.last_name, e.address FROM documents d`).
			WithArgs(input.CPF, input.Password).
			WillReturnError(expectedErr)
		result, err := repo.AuthenticationByCPF(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database error during authentication for CPF 12345678900")
		assert.Equal(t, readmodel.SignInReadModel{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should error from roles query", func(t *testing.T) {
		expectedErr := errors.New("roles query error")
		rows := pgxmock.NewRows([]string{"idusers", "first_name", "last_name", "address"}).
			AddRow(expectedIDUser, expectedFirstName, expectedLastName, expectedEmail)
		mock.ExpectQuery(`SELECT u.idusers, p.first_name, p.last_name, e.address FROM documents d`).
			WithArgs(input.CPF, input.Password).
			WillReturnRows(rows)
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(expectedIDUser).
			WillReturnError(expectedErr)
		result, err := repo.AuthenticationByCPF(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to get roles for user ID 123 during authentication")
		assert.ErrorContains(t, err, "roles query error")
		assert.Equal(t, readmodel.SignInReadModel{
			IDUser:    expectedIDUser,
			FirstName: expectedFirstName,
			LastName:  expectedLastName,
			Email:     expectedEmail,
		}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUser_GetRolesByAccountID(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	// db := &database.DB{Pool: mock}
	repo := NewUser(mock)
	idUser := 123
	t.Run("should return roles", func(t *testing.T) {
		expectedRoles := []string{"admin", "manager", "user"}
		rows := pgxmock.NewRows([]string{"name"}).
			AddRow("admin").
			AddRow("manager").
			AddRow("user")
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(idUser).
			WillReturnRows(rows)
		roles, err := repo.GetRolesByID(ctx, idUser)
		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should empty roles", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"name"})
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(idUser).
			WillReturnRows(rows)
		// roles, err := repo.GetPermissionsByRoleID(ctx, idUser)
		roles, err := repo.GetRolesByID(ctx, idUser)
		assert.NoError(t, err)
		assert.Nil(t, roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error", func(t *testing.T) {
		expectedErr := errors.New("query error")
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(idUser).
			WillReturnError(expectedErr)
		roles, err := repo.GetRolesByID(ctx, idUser)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database error while getting roles for user ID 123")
		// assert.Equal(t, expectedErr, err)
		assert.Nil(t, roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should error from database", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"name"}).
			AddRow("admin").
			AddRow(nil).
			RowError(1, errors.New("scan error"))
		mock.ExpectQuery(`SELECT r.name FROM roles r`).
			WithArgs(idUser).
			WillReturnRows(rows)
		roles, err := repo.GetRolesByID(ctx, idUser)
		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUser_GetPermissionsByRoleID(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	// db := &database.DB{Pool: mock}
	repo := NewUser(mock)
	t.Run("should return permissions", func(t *testing.T) {
		roleID := 1
		expectedPermissions := []readmodel.Permission{
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
		assert.ErrorContains(t, err, "database error while getting permissions for role ID 1")
		assert.ErrorContains(t, err, "query error")
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

func TestUser_FindEmailByCPF(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	// db := &database.DB{Pool: mock}
	repo := NewUser(mock)
	t.Run("should return email by cpf", func(t *testing.T) {
		cpf := "11144477735"
		expectedEmail := readmodel.GetEmailByCPFReadModel{Email: "email@example.com.br"}
		rows := pgxmock.NewRows([]string{"address"}).AddRow(expectedEmail.Email)
		mock.ExpectQuery(`SELECT e.address FROM emails e`).
			WithArgs(cpf).
			WillReturnRows(rows)
		email, err := repo.FindEmailByCPF(ctx, cpf)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmail, email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when email not found", func(t *testing.T) {
		cpf := "00000000000"
		mock.ExpectQuery(`SELECT e.address FROM emails e`).
			WithArgs(cpf).
			WillReturnError(pgx.ErrNoRows)
		email, err := repo.FindEmailByCPF(ctx, cpf)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "email not found for CPF 00000000000")
		assert.ErrorContains(t, err, "email not found")
		assert.Equal(t, "", email.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return database error", func(t *testing.T) {
		cpf := "12345678900"
		mockErr := errors.New("database error")
		mock.ExpectQuery(`SELECT e.address FROM emails e`).
			WithArgs(cpf).
			WillReturnError(mockErr)
		email, err := repo.FindEmailByCPF(ctx, cpf)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database error while finding email for CPF 12345678900")
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Equal(t, "", email.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUser_FindByEmail(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	// db := &database.DB{Pool: mock}
	repo := NewUser(mock)

	t.Run("should return an id user", func(t *testing.T) {
		email := "manager@email.com.br"
		expectedIDUser := 123
		mock.ExpectQuery(`SELECT idusers FROM users u`).
			WithArgs(email).
			WillReturnRows(pgxmock.NewRows([]string{"idusers"}).AddRow(expectedIDUser))
		idUser, err := repo.FindIDUserByEmail(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, expectedIDUser, idUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		email := "envalid@email.com.br"
		mock.ExpectQuery(`SELECT idusers FROM users u`).
			WithArgs(email).
			WillReturnError(pgx.ErrNoRows)
		_, err := repo.FindIDUserByEmail(ctx, email)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "user not found for email")
	})

	t.Run("should return database error", func(t *testing.T) {
		email := "email@email.com.br"
		expectedErr := errors.New("database error")
		mock.ExpectQuery(`SELECT idusers FROM users u`).
			WithArgs(email).
			WillReturnError(expectedErr)
		_, err := repo.FindIDUserByEmail(ctx, email)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database error while finding user for email")
		assert.ErrorContains(t, err, expectedErr.Error())
	})
}

func TestUser_UpdatePassword(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	// db := &database.DB{Pool: mock}
	repo := NewUser(mock)
	password := "@Dmin1234"
	cpf := "11144477735"

	t.Run("should successful", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users`).
			WithArgs(123, password, cpf).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		err := repo.UpdatePassword(ctx, 123, password, cpf)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error if cpf is not users", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users`).
			WithArgs(123, password, "98765432100").
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))
		err := repo.UpdatePassword(ctx, 123, password, "98765432100")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no rows affected while updating password for user ID 123")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error", func(t *testing.T) {
		mockErr := errors.New("database error")
		mock.ExpectExec(`UPDATE users`).
			WithArgs(123, "password", cpf).
			WillReturnError(mockErr)
		err := repo.UpdatePassword(ctx, 123, "password", cpf)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database error while updating password for user ID 123")
		assert.ErrorContains(t, err, mockErr.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
