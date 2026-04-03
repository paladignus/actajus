// Package postgres
package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	sharedpostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type CatalogQuery struct {
	db sharedpostgres.Executor
}

func NewCatalogQuery(db sharedpostgres.Executor) CatalogQuery {
	return CatalogQuery{db: db}
}

func (r CatalogQuery) ListUsers(ctx context.Context, filter repository.CatalogUserFilter) ([]readmodel.UserListItemReadModel, error) {
	limit := normalizeCatalogLimit(filter.Limit)
	query := "%" + strings.TrimSpace(strings.ToLower(filter.Query)) + "%"
	const statement = `
		SELECT
			u.idusers,
			TRIM(COALESCE(p.first_name, '') || ' ' || COALESCE(p.last_name, '')) AS full_name,
			COALESCE(e.address, '') AS email,
			COALESCE(u.is_blocked, FALSE) AS is_blocked,
			u.last_login_at
		FROM users u
		LEFT JOIN people p ON p.idpeople = u.idusers
		LEFT JOIN email_person ep ON ep.id_people = u.idusers
		LEFT JOIN emails e ON e.idemails = ep.id_emails AND e.is_primary = TRUE AND e.deleted_at IS NULL
		WHERE u.deleted_at IS NULL
		  AND ($1 = '%%'
		    OR LOWER(COALESCE(p.first_name, '')) LIKE $1
		    OR LOWER(COALESCE(p.last_name, '')) LIKE $1
		    OR LOWER(COALESCE(e.address, '')) LIKE $1)
		ORDER BY u.idusers DESC
		LIMIT $2;`
	rows, err := r.db.Query(ctx, statement, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]readmodel.UserListItemReadModel, 0, limit)
	for rows.Next() {
		var item readmodel.UserListItemReadModel
		var lastLoginAt *time.Time
		if err := rows.Scan(&item.ID, &item.FullName, &item.Email, &item.IsBlocked, &lastLoginAt); err != nil {
			return nil, err
		}
		item.LastLoginAt = lastLoginAt
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r CatalogQuery) ListPermissions(ctx context.Context, filter repository.CatalogPermissionFilter) ([]readmodel.PermissionListItemReadModel, error) {
	limit := normalizeCatalogLimit(filter.Limit)
	query := "%" + strings.TrimSpace(strings.ToLower(filter.Query)) + "%"
	const statement = `
		SELECT
			p.idpermissions,
			p.resource,
			p.action,
			p.description,
			p.created_at
		FROM permissions p
		WHERE $1 = '%%'
		   OR LOWER(p.resource) LIKE $1
		   OR LOWER(p.action) LIKE $1
		   OR LOWER(p.description) LIKE $1
		ORDER BY p.resource ASC, p.action ASC, p.idpermissions ASC
		LIMIT $2;`
	rows, err := r.db.Query(ctx, statement, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]readmodel.PermissionListItemReadModel, 0, limit)
	for rows.Next() {
		var item readmodel.PermissionListItemReadModel
		if err := rows.Scan(&item.ID, &item.Resource, &item.Action, &item.Description, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r CatalogQuery) ListRoles(ctx context.Context) ([]readmodel.RoleListItemReadModel, error) {
	const statement = `
		SELECT
			r.idroles,
			r.code,
			r.description,
			r.created_at
		FROM roles r
		ORDER BY r.code ASC, r.idroles ASC;`
	rows, err := r.db.Query(ctx, statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]readmodel.RoleListItemReadModel, 0, 16)
	for rows.Next() {
		var item readmodel.RoleListItemReadModel
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r CatalogQuery) ListRolesByUser(ctx context.Context, idUser int64) ([]readmodel.UserRoleAssignmentReadModel, error) {
	const statement = `
		SELECT
			ru.id_users,
			ru.id_roles,
			r.code,
			ru.assigned_by
		FROM role_user ru
		JOIN roles r ON r.idroles = ru.id_roles
		WHERE ru.id_users = $1
		ORDER BY r.code ASC, ru.id_roles ASC;`
	rows, err := r.db.Query(ctx, statement, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]readmodel.UserRoleAssignmentReadModel, 0, 8)
	for rows.Next() {
		var item readmodel.UserRoleAssignmentReadModel
		if err := rows.Scan(&item.IDUser, &item.IDRole, &item.RoleName, &item.AssignedBy); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r CatalogQuery) ListPermissionsByRole(ctx context.Context, idRole int16) ([]readmodel.RolePermissionAssignmentReadModel, error) {
	const statement = `
		SELECT
			pr.id_roles,
			pr.id_permissions,
			p.resource,
			p.action,
			p.description
		FROM permission_role pr
		JOIN permissions p ON p.idpermissions = pr.id_permissions
		WHERE pr.id_roles = $1
		ORDER BY p.resource ASC, p.action ASC, p.idpermissions ASC;`
	rows, err := r.db.Query(ctx, statement, idRole)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]readmodel.RolePermissionAssignmentReadModel, 0, 16)
	for rows.Next() {
		var item readmodel.RolePermissionAssignmentReadModel
		if err := rows.Scan(&item.IDRole, &item.IDPermission, &item.Resource, &item.Action, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func normalizeCatalogLimit(limit int) int {
	if limit <= 0 || limit > 200 {
		return 50
	}
	return limit
}

var _ repository.CatalogQueryRepository = CatalogQuery{}
