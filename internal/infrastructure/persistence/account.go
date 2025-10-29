// Package persistence
package persistence

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type Account struct {
	db *database.DB
}

func NewAccount(db *database.DB) Account {
	return Account{db}
}

func (a Account) FindUserAccountByCPF(ctx context.Context, cpf string) (authenticated dto.AuthenticateOutput, err error) {
	var IDAccount string
	sql := `SELECT p.idpeople, p.first_name, p.last_name, e.address, a.idaccounts FROM people p
	INNER JOIN documents d ON p.idpeople = d.id_people
	INNER JOIN accounts a ON p.idpeople = a.id_people
	INNER JOIN emails e ON p.idpeople = e.id_people
	WHERE p.deleted_at IS NULL AND a.deleted_at IS NULL AND d.cpf = $1;`
	if err = a.db.Pool.QueryRow(ctx, sql, cpf).
		Scan(
			&authenticated.IDUser,
			&authenticated.FirstName,
			&authenticated.LastName,
			&authenticated.Email,
			&IDAccount,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.AuthenticateOutput{}, exception.ErrUserNotFound
		}
		return dto.AuthenticateOutput{}, err
	}
	roles, err := a.GetRolesByAccountID(ctx, IDAccount)
	if err != nil {
		return dto.AuthenticateOutput{}, err
	}
	authenticated.Roles = roles
	return authenticated, nil
}

func (a Account) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
	var ok bool
	sql := `SELECT EXISTS (
      SELECT 1 FROM accounts a
      WHERE a.id_people = $1 AND a.password = crypt($2, password)
    ) AS valid`
	if err := a.db.Pool.QueryRow(ctx, sql, IDPeople, password).
		Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return exception.ErrInvalidCredentials
	}
	return nil
}

func (a Account) GetRolesByAccountID(ctx context.Context, idaccount string) (roles []string, err error) {
	sql := `SELECT r.name FROM roles r
		INNER JOIN account_role ar ON r.idroles = ar.id_roles
		WHERE ar.id_accounts = $1 ORDER BY r.name;`
	rows, err := a.db.Pool.Query(ctx, sql, idaccount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (a Account) GetPermissionsByRoleID(ctx context.Context, idrole int) (permissions []dto.Permission, err error) {
	sql := `SELECT p.resource, p.action FROM permissions p
		INNER JOIN role_permission rp ON p.idpermissions = rp.id_permissions
		WHERE rp.id_roles = $1 ORDER BY p.resource, p.action;`
	rows, err := a.db.Pool.Query(ctx, sql, idrole)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var permission dto.Permission
		if err := rows.Scan(&permission.Resource, &permission.Action); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (a Account) FindEmailByCPF(ctx context.Context, cpf string) (email string, err error) {
	sql := `SELECT e.address FROM emails e
		INNER JOIN people p ON e.id_people = p.idpeople
		INNER JOIN documents d ON d.id_people = p.idpeople
		WHERE e.deleted_at IS NULL AND p.deleted_at IS NULL AND d.cpf = $1;`
	if err = a.db.Pool.QueryRow(ctx, sql, cpf).Scan(&email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", exception.ErrEmailNotFound
		}
		return "", err
	}
	return email, nil
}
