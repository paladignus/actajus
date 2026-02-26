// Package postgres
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/module/identity/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type User struct {
	db postgres.PgxPool
}

func NewUser(db postgres.PgxPool) User {
	return User{db}
}

func (r User) FindByEmail(ctx context.Context, email vo.Email) (*domain.User, error) {
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
	err := r.db.QueryRow(ctx, query, email.Value()).Scan(
		&id, &passwordHash, &isBlocked,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return domain.NewUserBuilder().
		WithID(domain.IDUser(id)).
		WithPrimaryEmail(email).
		WithPasswordHash(domain.PasswordHash(passwordHash)).
		WithIsBlocked(isBlocked).
		Build()
}

func (r User) FindByID(ctx context.Context, id domain.IDUser) (*domain.User, error) {
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
	err := r.db.QueryRow(ctx, query, id.Value()).Scan(
		&passwordHash, &isBlocked, &email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return domain.NewUserBuilder().
		WithID(domain.IDUser(id)).
		WithPrimaryEmail(vo.Email(email)).
		WithPasswordHash(domain.PasswordHash(passwordHash)).
		WithIsBlocked(isBlocked).
		Build()
}

func (r User) UpdatePasswordHash(ctx context.Context, idUser domain.IDUser, passwordHash string) error {
	const query = `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE idusers = $2 AND deleted_at IS NULL`
	ct, err := r.db.Exec(ctx, query, passwordHash, idUser.Value())
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
