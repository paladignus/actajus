// Package postgres
package postgres

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/model"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type PasswordReset struct {
	db postgres.Executor
}

func NewPasswordReset(db postgres.Executor) *PasswordReset {
	return &PasswordReset{db}
}

func (r PasswordReset) Create(ctx context.Context, input *model.PasswordResetTokenCreate) (int64, error) {
	const query = `INSERT INTO
  password_resets (id_users, token_hash, expires_at, created_at)
  VALUES ($1, $2, $3, $4) RETURNING idpassword_resets;`
	var id int64
	if err := r.db.QueryRow(ctx, query,
		input.IDUser,
		input.Hash[:],
		input.ExpiresAt,
		input.CreatedAt,
	).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r PasswordReset) GetByID(ctx context.Context, id int64) (*model.PasswordResetToken, error) {
	const query = `SELECT
	id_users, token_hash, expires_at, used_at
	FROM password_resets
	WHERE idpassword_resets = $1 LIMIT 1;`
	var (
		uid       int64
		hashBytes []byte
		expiresAt time.Time
		usedAt    *time.Time
	)
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&uid, &hashBytes, &expiresAt, &usedAt,
	); err != nil {
		return nil, err
	}
	hash, err := scanSHA256(hashBytes)
	if err != nil {
		return nil, err
	}
	return &model.PasswordResetToken{
		IDUser:    uid,
		Hash:      hash,
		ExpiresAt: expiresAt,
		UsedAt:    usedAt,
	}, nil
}

func (r PasswordReset) MarkUsed(ctx context.Context, id int64, usedAt time.Time) error {
	const query = `UPDATE password_resets
  SET used_at = $2
  WHERE idpassword_resets = $1 AND used_at IS NULL;`
	_, err := r.db.Exec(ctx, query, id, usedAt)
	return err
}

func (r PasswordReset) RevokeAllByUser(ctx context.Context, uid int64, now time.Time) error {
	const query = `UPDATE password_resets
  SET used_at = $1
  WHERE id_users = $2 AND revoked_at IS NULL;`
	_, err := r.db.Exec(ctx, query, now, uid)
	return err
}

func (r PasswordReset) GetActiveByUser(ctx context.Context, uid int64) (*model.PasswordResetToken, error) {
	const query = `SELECT
	  idpassword_resets, token_hash, expires_at, used_at
	FROM password_resets
	WHERE id_users = $1 AND used_at IS NULL AND expires_at > NOW() LIMIT 1;`
	var (
		prid      int64
		hashBytes []byte
		expiresAt time.Time
		usedAt    *time.Time
	)
	if err := r.db.QueryRow(ctx, query, uid).Scan(
		&prid, &hashBytes, &expiresAt, &usedAt,
	); err != nil {
		return nil, err
	}
	hash, err := scanSHA256(hashBytes)
	if err != nil {
		return nil, err
	}
	return &model.PasswordResetToken{
		ID:        prid,
		Hash:      hash,
		ExpiresAt: expiresAt,
		UsedAt:    usedAt,
	}, nil
}
