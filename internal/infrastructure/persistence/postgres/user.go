// Package postgres
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type User struct {
	db *database.DB
}

func NewUser(db *database.DB) User {
	return User{db}
}

func (u User) AuthenticationByCPF(ctx context.Context, input dto.SignInInput) (user dto.SignInOutput, err error) {
	sql := `
	SELECT
    u.idusers,
    p.first_name,
    p.last_name,
    e.address
	FROM documents d
	JOIN people p ON p.idpeople = d.id_people AND p.deleted_at IS NULL
	JOIN users u ON u.idusers = p.idpeople AND u.deleted_at IS NULL
	LEFT JOIN emails e ON e.id_people = p.idpeople AND e.deleted_at IS NULL
	WHERE cpf = $1 AND u.password = crypt($2, password);`
	if err = u.db.Pool.QueryRow(ctx, sql, input.CPF, input.Password).
		Scan(
			&user.IDUser,
			&user.FirstName,
			&user.LastName,
			&user.Email,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("authentication failed for CPF %s: %w", input.CPF, exception.ErrInvalidCredentials)
		}
		return user, fmt.Errorf("database error during authentication for CPF %s: %w", input.CPF, err)
	}
	roles, err := u.GetRolesByID(ctx, user.IDUser)
	if err != nil {
		return user, fmt.Errorf("failed to get roles for user ID %d during authentication: %w", user.IDUser, err)
	}
	user.Roles = roles
	return user, nil
}

func (u User) GetRolesByID(ctx context.Context, idUser int) (roles []string, err error) {
	sql := `SELECT r.name FROM roles r
		INNER JOIN role_user ru ON r.idroles = ru.id_roles
		WHERE ru.id_users = $1 ORDER BY r.name;`
	rows, err := u.db.Pool.Query(ctx, sql, idUser)
	if err != nil {
		return nil, fmt.Errorf("database error while getting roles for user ID %d: %w", idUser, err)
	}
	defer rows.Close()
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("failed to scan role for user ID %d: %w", idUser, err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (u User) GetPermissionsByRoleID(ctx context.Context, idrole int) (permissions []dto.Permission, err error) {
	sql := `SELECT p.resource, p.action FROM permissions p
		INNER JOIN role_permission rp ON p.idpermissions = rp.id_permissions
		WHERE rp.id_roles = $1 ORDER BY p.resource, p.action;`
	rows, err := u.db.Pool.Query(ctx, sql, idrole)
	if err != nil {
		return nil, fmt.Errorf("database error while getting permissions for role ID %d: %w", idrole, err)
	}
	defer rows.Close()
	for rows.Next() {
		var permission dto.Permission
		if err := rows.Scan(&permission.Resource, &permission.Action); err != nil {
			return nil, fmt.Errorf("failed to scan permission for role ID %d: %w", idrole, err)
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (u User) FindEmailByCPF(ctx context.Context, cpf string) (output dto.GetEmailByCPFOutput, err error) {
	sql := `
		SELECT e.address FROM emails e
		LEFT JOIN documents d ON d.id_people = e.id_people
		JOIN users u ON u.idusers = e.id_people AND u.deleted_at IS NULL
		WHERE d.cpf = $1 AND e.deleted_at IS NULL;`
	if err = u.db.Pool.QueryRow(ctx, sql, cpf).Scan(&output.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return output, fmt.Errorf("email not found for CPF %s: %w", cpf, exception.ErrEmailNotFound)
		}
		return output, fmt.Errorf("database error while finding email for CPF %s: %w", cpf, err)
	}
	return output, nil
}

func (u User) FindIDUserByEmail(ctx context.Context, email string) (idUser int, err error) {
	sql := `
		SELECT idusers FROM users u
		LEFT JOIN emails e ON e.id_people = u.idusers AND e.deleted_at IS NULL
		WHERE e.address = $1 AND u.deleted_at IS NULL`
	if err := u.db.Pool.QueryRow(ctx, sql, email).Scan(&idUser); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return idUser, fmt.Errorf("user not found for email %s: %w", email, exception.ErrEmailNotFound)
		}
		return idUser, fmt.Errorf("database error while finding user for email %s: %w", email, err)
	}
	return idUser, nil
}

func (u User) UpdatePassword(ctx context.Context, idUser int, password, cpf string) error {
	sql := `
		UPDATE users u
    SET password = crypt($2, gen_salt('bf')), updated_at = now()
		FROM documents d
    WHERE idusers = $1 AND u.idusers = d.id_people AND d.cpf = $3 AND deleted_at IS NULL;`
	ct, err := u.db.Pool.Exec(ctx, sql, idUser, password, cpf)
	if err != nil {
		return fmt.Errorf("database error while updating password for user ID %d: %w", idUser, err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected while updating password for user ID %d: %w", idUser, exception.ErrInvalidCredentials)
	}
	return nil
}
