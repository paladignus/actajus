// Package postgres
package postgres

import (
	"context"

	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type RoleUserPair struct {
	IDRole int16
	IDUser int64
}

func ListAllRoleUsers(ctx context.Context, db postgres.Executor) ([]RoleUserPair, error) {
	const query = `SELECT
	id_roles, id_users
	FROM role_user;`
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]RoleUserPair, 0, 128)
	for rows.Next() {
		var p RoleUserPair
		if err := rows.Scan(&p.IDRole, &p.IDUser); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
