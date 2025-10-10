// Package persistence
package persistence

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
)

type Account struct {
	db *postgres.DB
}

func NewAccount(db *postgres.DB) Account {
	return Account{db}
}

func (a Account) FindUserAccountByCPF(ctx context.Context, cpf string) (authenticated dto.AuthenticatedOutput, err error) {
	var idaccount string
	sql := `SELECT p.idpeople, p.first_name, p.last_name, e.address, a.idaccounts FROM people p
	INNER JOIN documents d ON p.idpeople = d.id_people
	INNER JOIN accounts a ON p.idpeople = a.id_people
	INNER JOIN emails e ON p.idpeople = e.id_people
	WHERE p.deleted_at IS NULL AND a.deleted_at IS NULL AND d.cpf = $1;`
	err = a.db.Pool.QueryRow(ctx, sql, cpf).
		Scan(&authenticated.IDUser, &authenticated.FirstName, &authenticated.LastName, &authenticated.Email, &idaccount)
	if errors.Is(err, pgx.ErrNoRows) {
		return authenticated, domain.ErrPersonNotFound
	}
	return authenticated, err
}

func (a Account) ValidatePassword(ctx context.Context, IDPeople, password string) (ok bool) {
	sql := `SELECT EXISTS (
      SELECT 1 FROM accounts a
      WHERE a.deleted_at IS NULL AND a.id_people = $1 AND a.password = crypt($2, password)
    ) AS valid`
	a.db.Pool.QueryRow(ctx, sql, IDPeople, password).Scan(&ok)
	return ok
}

// type Authenticate struct {
// 	db *postgres.DB
// }
//
// func NewAuthenticate(db *postgres.DB) Authenticate {
// 	return Authenticate{db}
// }
//
// func (a Authenticate) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
// 	autenticated := dto.AuthenticatedOutput{}
// 	var idaccount string
// 	sql := `SELECT p.idpeople, p.first_name, p.last_name, e.address, a.idaccounts FROM people p
// 		INNER JOIN documents d ON p.idpeople = d.id_people
// 		INNER JOIN accounts a ON p.idpeople = a.id_people
// 		INNER JOIN emails e ON p.idpeople = e.id_people
// 		WHERE d.cpf = $1 AND a.password = crypt($2, password);`
// 	err := a.db.Pool.QueryRow(ctx, sql, cpf, password).
// 		Scan(&autenticated.UserID, &autenticated.FirstName, &autenticated.LastName, &autenticated.Email, &idaccount)
// 	if errors.Is(err, pgx.ErrNoRows) {
// 		return autenticated, domain.ErrUserNotFound
// 	}
// 	roles, err := a.GetAccountRoles(ctx, idaccount)
// 	if err != nil {
// 		return autenticated, domain.ErrUserNotFound
// 	}
// 	autenticated.Roles = append(autenticated.Roles, roles...)
// 	return autenticated, err
// }
//
// func (a Authenticate) GetAccountRoles(ctx context.Context, idaccount string) (roles []string, err error) {
// 	sql := `SELECT r.name FROM roles r
// 		INNER JOIN account_role ar ON r.idroles = ar.id_roles
// 		WHERE ar.id_accounts = $1 ORDER BY r.name;`
// 	rows, err := a.db.Pool.Query(ctx, sql, idaccount)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var role string
// 		if err := rows.Scan(&role); err != nil {
// 			return nil, err
// 		}
// 		roles = append(roles, role)
// 	}
// 	return roles, nil
// }
//
// func (a Authenticate) GetRolePermissions(ctx context.Context, idrole int) (permissions []dto.Permission, err error) {
// 	sql := `SELECT p.resource, p.action FROM permissions p
// 		INNER JOIN role_permission rp ON p.idpermissions = rp.id_permissions
// 		WHERE rp.id_roles = $1 ORDER BY p.resource, p.action;`
// 	rows, err := a.db.Pool.Query(ctx, sql, idrole)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var permission dto.Permission
// 		if err := rows.Scan(&permission.Resource, &permission.Action); err != nil {
// 			return nil, err
// 		}
// 		permissions = append(permissions, permission)
// 	}
// 	return permissions, nil
// }
