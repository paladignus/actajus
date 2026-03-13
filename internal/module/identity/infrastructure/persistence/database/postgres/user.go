// Package postgres
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/module/identity/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type User struct {
	db postgres.Executor
}

func NewUser(db postgres.Executor) *User {
	return &User{db}
}

func (r User) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT
			u.idusers, u.password_hash, u.is_blocked 
		FROM users u
		LEFT JOIN email_person ep ON u.idusers = ep.id_people
		JOIN emails e ON e.idemails = ep.id_emails AND e.is_primary = TRUE AND e.deleted_at IS NULL
		WHERE e.address = $1 AND u.deleted_at IS NULL LIMIT 1;`
	var (
		id           int64
		passwordHash string
		isBlocked    bool
	)
	err := r.db.QueryRow(ctx, query, email).Scan(
		&id, &passwordHash, &isBlocked,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return domain.NewUserBuilder().
		WithID(id).
		WithPrimaryEmail(email).
		WithPasswordHash(passwordHash).
		WithIsBlocked(isBlocked).
		Build()
}

func (r User) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	const query = `
		SELECT
			u.password_hash, u.is_blocked, e.address 
		FROM users u
		LEFT JOIN email_person ep ON u.idusers = ep.id_people
		JOIN emails e ON e.idemails = ep.id_emails AND e.is_primary = TRUE AND e.deleted_at IS NULL
		WHERE u.idusers = $1 AND u.deleted_at IS NULL LIMIT 1;`
	var (
		passwordHash string
		isBlocked    bool
		email        string
	)
	err := r.db.QueryRow(ctx, query, id).Scan(
		&passwordHash, &isBlocked, &email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return domain.NewUserBuilder().
		WithID(id).
		WithPrimaryEmail(email).
		WithPasswordHash(passwordHash).
		WithIsBlocked(isBlocked).
		Build()
}

func (r User) UpdatePasswordHash(ctx context.Context, id int64, passwordHash string) error {
	const query = `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE idusers = $2 AND deleted_at IS NULL`
	ct, err := r.db.Exec(ctx, query, passwordHash, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
