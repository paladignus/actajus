// Package postgres
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/module/identity/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type Session struct {
	db postgres.PgxPool
}

func NewSession(db postgres.PgxPool) Session {
	return Session{db}
}

func (r Session) Create(ctx context.Context, s *domain.Session) error {
	const query = `
		INSERT INTO sessions (
			id_users, refresh_token_hash, expires_at, revoked_at,
		  rotated_at, ip, user_agent, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9
		) RETURNING idsessions;`
	var id int64
	h := s.RefreshHash()
	if err := r.db.QueryRow(ctx, query,
		s.IDUser().Value(),
		h[:],
		s.ExpiresAt(),
		s.RevokedAt(),
		s.RotatedAt(),
		s.IP(),
		s.UserAgent(),
		s.CreatedAt(),
		s.UpdatedAt(),
	).Scan(&id); err != nil {
		return err
	}
	return s.SetID(id)
}

func (r Session) GetByID(ctx context.Context, id int64) (*domain.Session, error) {
	const q = `SELECT
	id_users, refresh_token_hash, expires_at, revoked_at, rotated_at,
	COALESCE(ip, ''), COALESCE(user_agent, ''), created_at, updated_at
	FROM sessions
	WHERE idsessions = $1 LIMIT 1;`
	var (
		uid       int64
		hashBytes []byte
		expiresAt time.Time
		revokedAt *time.Time
		rotatedAt *time.Time
		ip        string
		userAgent string
		createdAt time.Time
		updatedAt time.Time
	)
	if err := r.db.QueryRow(ctx, q, id).Scan(
		&uid, &hashBytes, &expiresAt, &revokedAt, &rotatedAt,
		&ip, &userAgent, &createdAt, &updatedAt,
	); err != nil {
		return nil, err
	}
	hash, err := scanSHA256(hashBytes)
	if err != nil {
		return nil, err
	}
	return domain.NewSessionBuilder().
		WithID(id).
		WithIDUser(uid).
		WithRefreshHash(hash).
		WithExpiresAt(expiresAt).
		WithRevokedAt(revokedAt).
		WithRotatedAt(rotatedAt).
		WithIP(ip).
		WithUserAgent(userAgent).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()
}

func (r Session) RotateRefreshToken(ctx context.Context, id int64, hash [32]byte, expiresAt time.Time) error {
	const q = `UPDATE sessions
	SET refresh_token_hash = $2, expires_at = $3, rotated_at = $4, updated_at = $4
	WHERE idsessions = $1 AND revoked_at IS NULL;`
	now := time.Now()
	if _, err := r.db.Exec(ctx, q,
		id,
		hash[:],
		expiresAt,
		now,
	); err != nil {
		return err
	}
	return nil
}

func (r Session) Revoke(ctx context.Context, id int64) error {
	const q = `UPDATE sessions
	SET revoked_at = $2, updated_at = $2
	WHERE idsessions = $1 AND revoked_at IS NULL;`
	now := time.Now()
	if _, err := r.db.Exec(ctx, q,
		id, now); err != nil {
		return err
	}
	return nil
}

func (r Session) RevokeAllByUser(ctx context.Context, uid int64) error {
	const q = `UPDATE sessions
	SET revoked_at = $2, updated_at = $2
	WHERE id_users = $1 AND revoked_at IS NULL;`
	now := time.Now()
	if _, err := r.db.Exec(ctx, q,
		uid, now); err != nil {
		return err
	}
	return nil
}

func (r Session) CountActiveByUser(ctx context.Context, uid int64) (int, error) {
	const q = `SELECT COUNT(*)
	FROM sessions
	WHERE id_users = $1 AND revoked_at IS NULL AND expires_at > NOW();`
	var n int
	if err := r.db.QueryRow(ctx, q,
		uid).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func scanSHA256(b []byte) ([32]byte, error) {
	var out [32]byte
	if len(b) != 32 {
		return out, fmt.Errorf("invalid sha256 length: %d", len(b))
	}
	copy(out[:], b)
	return out, nil
}

func (r *Session) IsActive(ctx context.Context, id int64, uid int64, now time.Time) (bool, error) {
	const query = `SELECT 1
	FROM sessions
	WHERE idsessions = $1 AND user_id = $2 AND revoked_at IS NULL AND expires_at > $3 LIMIT 1;`
	var one int
	err := r.db.QueryRow(ctx, query, id, uid, now).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Session) RotateRefreshTokenAtomic(
	ctx context.Context,
	id int64,
	oldHash [32]byte,
	hash [32]byte,
	expiresAt time.Time,
	now time.Time,
) (bool, error) {
	const query = `UPDATE sessions
	SET refresh_token_hash = $1, expires_at = $2, rotated_at = $3, updated_at = $3
	WHERE idsessions = $4 AND revoked_at IS NULL AND refresh_token_hash = $5
	RETURNING idsessions;`
	old := oldHash
	nw := hash
	var outID int64
	if err := r.db.QueryRow(ctx, query,
		nw[:],
		expiresAt,
		now,
		id,
		old[:]).Scan(&outID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// func nullableText(s string) any {
// 	if s == "" {
// 		return nil
// 	}
// 	return s
// }
