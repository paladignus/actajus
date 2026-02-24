// Package postgres
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	identityRepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	identity "github.com/paladignus/actajus/internal/module/identity/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

var _ identityRepo.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) FindByEmail(ctx context.Context, email vo.Email) (*identity.User, error) {
	const query = `
SELECT
	u.id,
	u.password_hash,
	u.is_blocked,
	u.created_at,
	u.updated_at,
FROM users u
JOIN emails e
	ON e.user_id = u.id
	AND e.is_primary = TRUE
	AND e.deleted_at IS NULL
WHERE
	e.address = $1
	AND u.deleted_at IS NULL
LIMIT 1;`

	var (
		id           int64
		passwordHash string
		isBlocked    bool
		createdAt    time.Time
		updatedAt    time.Time
		deletedAt    *time.Time
	)

	err := r.pool.QueryRow(ctx, query, email.Value()).Scan(
		&id,
		&passwordHash,
		&isBlocked,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return identity.NewUserFromPersistence(
		identity.UserID(id),
		identity.PasswordHash(passwordHash),
		isBlocked,
		createdAt,
		updatedAt,
		deletedAt,
	), nil
}

func (r *UserRepository) FindByID(ctx context.Context, uid identity.UserID) (*identity.User, error) {
	const q = `
SELECT
	u.id,
	u.password_hash,
	u.is_blocked,
	u.created_at,
	u.updated_at,
	u.deleted_at
FROM users u
WHERE
	u.id = $1
	AND u.deleted_at IS NULL
LIMIT 1;`

	var (
		id           int64
		passwordHash string
		isBlocked    bool
		createdAt    time.Time
		updatedAt    time.Time
		deletedAt    *time.Time
	)

	err := r.pool.QueryRow(ctx, q, uid.Value()).Scan(
		&id,
		&passwordHash,
		&isBlocked,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return identity.NewUserFromPersistence(
		identity.UserID(id),
		identity.PasswordHash(passwordHash),
		isBlocked,
		createdAt,
		updatedAt,
		deletedAt,
	), nil
}

func (r *UserRepository) UpdatePasswordHash(ctx context.Context, userID identity.UserID, passwordHash string) error {
	const q = `
UPDATE users
SET password_hash = $2,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;`

	ct, err := r.pool.Exec(ctx, q, userID.Value(), passwordHash)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
