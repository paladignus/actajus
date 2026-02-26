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
	return s.SetID(domain.IDSession(id))
}

func (r Session) GetByID(ctx context.Context, sid domain.IDSession) (*domain.Session, error) {
	const q = `SELECT
	id_users, refresh_token_hash, expires_at, revoked_at, rotated_at,
	COALESCE(ip, ''), COALESCE(user_agent, ''), created_at, updated_at
	FROM sessions
	WHERE idsessions = $1 LIMIT 1;`
	var (
		idUser    int64
		hashBytes []byte
		expiresAt time.Time
		revokedAt *time.Time
		rotatedAt *time.Time
		ip        string
		userAgent string
		createdAt time.Time
		updatedAt time.Time
	)
	if err := r.db.QueryRow(ctx, q, sid.Value()).Scan(
		&idUser, &hashBytes, &expiresAt, &revokedAt, &rotatedAt,
		&ip, &userAgent, &createdAt, &updatedAt,
	); err != nil {
		return nil, err
	}
	hash, err := scanSHA256(hashBytes)
	if err != nil {
		return nil, err
	}
	return domain.NewSessionBuilder().
		WithID(sid).
		WithIDUser(domain.IDUser(idUser)).
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

func (r Session) RotateRefreshToken(ctx context.Context, sid domain.IDSession, newHash [32]byte, newExpiresAt time.Time) error {
	const q = `UPDATE sessions
	SET refresh_token_hash = $2, expires_at = $3, rotated_at = $4, updated_at = $4
	WHERE idsessions = $1 AND revoked_at IS NULL;`
	now := time.Now()
	if _, err := r.db.Exec(ctx, q,
		sid.Value(),
		newHash[:],
		newExpiresAt,
		now,
	); err != nil {
		return err
	}
	return nil
}

func (r Session) Revoke(ctx context.Context, sid domain.IDSession) error {
	const q = `UPDATE sessions
	SET revoked_at = $2, updated_at = $2
	WHERE idsessions = $1 AND revoked_at IS NULL;`
	now := time.Now()
	if _, err := r.db.Exec(ctx, q,
		sid.Value(), now); err != nil {
		return err
	}
	return nil
}

func (r Session) RevokeAllByUser(ctx context.Context, idUser domain.IDUser) error {
	const q = `UPDATE sessions
	SET revoked_at = $2, updated_at = $2
	WHERE id_users = $1 AND revoked_at IS NULL;`
	now := time.Now()
	if _, err := r.db.Exec(ctx, q,
		idUser.Value(), now); err != nil {
		return err
	}
	return nil
}

func (r Session) CountActiveByUser(ctx context.Context, idUser domain.IDUser) (int, error) {
	const q = `SELECT COUNT(*)
	FROM sessions
	WHERE id_users = $1 AND revoked_at IS NULL AND expires_at > NOW();`
	var n int
	if err := r.db.QueryRow(ctx, q,
		idUser.Value()).Scan(&n); err != nil {
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

func (r *Session) IsActive(ctx context.Context, sid domain.IDSession, idUser domain.IDUser, now time.Time) (bool, error) {
	const query = `SELECT 1
	FROM sessions
	WHERE idsessions = $1 AND user_id = $2 AND revoked_at IS NULL AND expires_at > $3 LIMIT 1;`
	var one int
	err := r.db.QueryRow(ctx, query, sid.Value(), idUser.Value(), now).Scan(&one)
	if err != nil {
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
