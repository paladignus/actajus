// Package postgres
package postgres

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type Authorization struct {
	db postgres.PgxPool
}

func NewAuthorization(db postgres.PgxPool) Authorization {
	return Authorization{db}
}

func (r Authorization) ListPermissionsByUser(ctx context.Context, idUser domain.IDUser) ([]string, error) {
	const query = `SELECT
		DISTINCT (p.resource || ':' || p.action) AS perm
	FROM role_user ru
	JOIN permission_role pr ON pr.id_roles = ru.id_roles
	JOIN permissions p ON p.idpermissions = pr.id_permissions
	WHERE ru.id_users = $1;`
	rows, err := r.db.Query(ctx, query, idUser.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	perms := make([]string, 0, 32)
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return perms, nil
}
