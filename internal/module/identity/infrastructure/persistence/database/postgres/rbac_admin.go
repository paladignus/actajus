// Package postgres
package postgres

import (
	"context"

	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
)

type RoleUserAdminRepository struct {
	db postgres.PgxPool
}

func NewRoleUserAdminRepository(db postgres.PgxPool) RoleUserAdminRepository {
	return RoleUserAdminRepository{db}
}

func (r RoleUserAdminRepository) AssignRole(ctx context.Context, uid int64, rid int16, assignedBy int64) error {
	const query = `INSERT INTO role_user
	(id_users, id_roles, assigned_by) VALUES ($1, $2, $3)
	ON CONFLICT (id_roles, id_users) DO NOTHING;`
	_, err := r.db.Exec(ctx, query, uid, rid, assignedBy)
	return err
}

func (r RoleUserAdminRepository) RemoveRole(ctx context.Context, uid int64, rid int16) error {
	const query = `DELETE FROM role_user
	WHERE id_users = $1 AND id_roles = $2;`
	_, err := r.db.Exec(ctx, query, uid, rid)
	return err
}

func (r RoleUserAdminRepository) ListUserIDsByRole(ctx context.Context, rid int16) ([]int64, error) {
	const query = `SELECT id_users
	FROM role_user WHERE id_roles = $1;`
	rows, err := r.db.Query(ctx, query, rid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]int64, 0, 64)
	for rows.Next() {
		var uid int64
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		out = append(out, int64(uid))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

type PermissionRoleAdminRepository struct {
	db postgres.PgxPool
}

func NewPermissionRoleAdminRepository(db postgres.PgxPool) PermissionRoleAdminRepository {
	return PermissionRoleAdminRepository{db}
}

func (r PermissionRoleAdminRepository) GrantPermission(ctx context.Context, rid int16, pid int16) error {
	const query = `INSERT INTO permission_role
	(id_roles, id_permissions) VALUES ($1, $2)
	ON CONFLICT (id_roles, id_permissions) DO NOTHING;`
	_, err := r.db.Exec(ctx, query, rid, pid)
	return err
}

func (r PermissionRoleAdminRepository) RevokePermission(ctx context.Context, rid int16, pid int16) error {
	const query = `DELETE FROM permission_role
	WHERE id_roles = $1 AND id_permissions = $2;`
	_, err := r.db.Exec(ctx, query, rid, pid)
	return err
}
