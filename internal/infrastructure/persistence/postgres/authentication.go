// Package postgres
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type Authentication struct {
	db *database.DB
}

func NewAuthentication(db *database.DB) Authentication {
	return Authentication{db}
}

func (a Authentication) SignIn(ctx context.Context, cpf string) (output dto.SignInOutput, err error) {
	var IDAccount string
	sql := `SELECT DISTINCT p.idpeople, p.first_name, p.last_name, e.address, a.idaccounts
		FROM people p
		INNER JOIN documents d ON p.idpeople = d.id_people
		INNER JOIN accounts a ON p.idpeople = a.id_people AND a.deleted_at IS NULL
		INNER JOIN emails e ON p.idpeople = e.id_people AND e.deleted_at IS NULL
		WHERE p.deleted_at IS NULL AND d.cpf = $1;`
	if err = a.db.Pool.QueryRow(ctx, sql, cpf).
		Scan(
			&output.IDUser,
			&output.FirstName,
			&output.LastName,
			&output.Email,
			&IDAccount,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.SignInOutput{}, exception.ErrUserNotFound
		}
		return dto.SignInOutput{}, err
	}
	roles, err := a.GetRolesByAccountID(ctx, IDAccount)
	if err != nil {
		return dto.SignInOutput{}, err
	}
	output.Roles = roles
	return output, nil
}

func (a Authentication) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
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

func (a Authentication) GetRolesByAccountID(ctx context.Context, idaccount string) (roles []string, err error) {
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

func (a Authentication) GetPermissionsByRoleID(ctx context.Context, idrole int) (permissions []dto.Permission, err error) {
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

func (a Authentication) FindEmailByCPF(ctx context.Context, cpf string) (output dto.GetEmailByCPFOutput, err error) {
	sql := `SELECT e.address FROM emails e
		INNER JOIN people p ON e.id_people = p.idpeople
		INNER JOIN accounts a ON p.idpeople = a.id_people
		INNER JOIN documents d ON d.id_people = p.idpeople
		WHERE a.deleted_at IS NULL AND e.deleted_at IS NULL AND p.deleted_at IS NULL AND d.cpf = $1;`
	if err = a.db.Pool.QueryRow(ctx, sql, cpf).Scan(&output.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.GetEmailByCPFOutput{}, exception.ErrCPFNotFound
		}
		return dto.GetEmailByCPFOutput{}, err
	}
	return output, nil
}

func (a Authentication) AccountIsActive(ctx context.Context, email string) (IDUser string, err error) {
	sql := `SELECT e.id_people FROM emails e
		INNER JOIN people p ON e.id_people = p.idpeople AND p.deleted_at IS NULL
		INNER JOIN accounts a ON p.idpeople = a.id_people AND a.deleted_at IS NULL
		WHERE e.deleted_at IS NULL AND e.address = $1;`
	if err := a.db.Pool.QueryRow(ctx, sql, email).Scan(&IDUser); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", exception.ErrEmailNotFound
		}
		return "", err
	}
	return IDUser, nil
}

func (a Authentication) InvalidAllTokensByIDUser(ctx context.Context, IDUser string) (err error) {
	sql := `UPDATE password_reset SET used_at = now() WHERE id_people = $1;`
	_, err = a.db.Pool.Exec(ctx, sql, IDUser)
	return err
}

func (a Authentication) CreateRecoverPassword(ctx context.Context, IDUser, token string) (err error) {
	sql := `INSERT INTO password_reset (id_people, token, expires_at) VALUES ($1, $2, $3);`
	_, err = a.db.Pool.Exec(ctx, sql, IDUser, token, time.Now().Add(time.Minute*30))
	return err
}

// FindByEmail(ctx context.Context, email string) (entity.Authentication, error)
// 	FindById(ctx context.Context, idAuthentication int) (entity.Authentication, error)
// 	UpdatePassword(ctx context.Context, idAuthentication int, hashedPassword string) error

func (a Authentication) FindByEmail(ctx context.Context, email string) (auth entity.Authentication, err error) {
	sql := `SELECT idaccounts FROM accounts WHERE id_people = (SELECT id_people FROM emails WHERE address = $1);`
	if err := a.db.Pool.QueryRow(ctx, sql, email).Scan(&email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Authentication{}, exception.ErrEmailNotFound
		}
		return entity.Authentication{}, err
	}
	return entity.Authentication{}, nil
}

func (a Authentication) FindById(ctx context.Context, idAuthentication int) (auth entity.Authentication, err error) {
	return entity.Authentication{}, nil
}

func (Authentication) UpdatePassword(ctx context.Context, idAuthentication int, hashedPassword string) error {
	return nil
}
