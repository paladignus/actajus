// Package postgres
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	sharedpostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type RoleCatalogAdminRepository struct {
	db sharedpostgres.Executor
}

func NewRoleCatalogAdminRepository(db sharedpostgres.Executor) RoleCatalogAdminRepository {
	return RoleCatalogAdminRepository{db: db}
}

func (r RoleCatalogAdminRepository) Create(ctx context.Context, name, description string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO roles (name, description) VALUES ($1, $2);`, name, description)
	return err
}

func (r RoleCatalogAdminRepository) Update(ctx context.Context, id int16, name, description string) error {
	ct, err := r.db.Exec(ctx, `UPDATE roles SET name = $1, description = $2 WHERE idroles = $3;`, name, description, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r RoleCatalogAdminRepository) Delete(ctx context.Context, id int16) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM roles WHERE idroles = $1;`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

type PermissionCatalogAdminRepository struct {
	db sharedpostgres.Executor
}

func NewPermissionCatalogAdminRepository(db sharedpostgres.Executor) PermissionCatalogAdminRepository {
	return PermissionCatalogAdminRepository{db: db}
}

func (r PermissionCatalogAdminRepository) Create(ctx context.Context, resource, action, description string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO permissions (resource, action, description) VALUES ($1, $2, $3);`, resource, action, description)
	return err
}

func (r PermissionCatalogAdminRepository) Update(ctx context.Context, id int16, resource, action, description string) error {
	ct, err := r.db.Exec(ctx, `UPDATE permissions SET resource = $1, action = $2, description = $3 WHERE idpermissions = $4;`, resource, action, description, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r PermissionCatalogAdminRepository) Delete(ctx context.Context, id int16) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM permissions WHERE idpermissions = $1;`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
