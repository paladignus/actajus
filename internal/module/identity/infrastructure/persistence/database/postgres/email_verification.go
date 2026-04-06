// Package postgres
package postgres

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type EmailVerification struct {
	db postgres.Executor
}

func NewEmailVerification(db postgres.Executor) *EmailVerification {
	return &EmailVerification{db: db}
}

func (r EmailVerification) Create(ctx context.Context, input *mapper.EmailVerificationTokenCreate) (int64, error) {
	const query = `INSERT INTO
	email_verifications (id_users, id_emails, token_hash, expires_at, created_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING idemail_verifications;`
	var id int64
	if err := r.db.QueryRow(ctx, query,
		input.IDUser,
		input.IDEmail,
		input.Hash[:],
		input.ExpiresAt,
		input.CreatedAt,
	).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r EmailVerification) GetByID(ctx context.Context, id int64) (*mapper.EmailVerificationToken, error) {
	const query = `SELECT
		id_users, id_emails, token_hash, expires_at, used_at, created_at
	FROM email_verifications
	WHERE idemail_verifications = $1 LIMIT 1;`
	var (
		idUser    int64
		idEmail   int64
		hashBytes []byte
		expiresAt time.Time
		usedAt    *time.Time
		createdAt time.Time
	)
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&idUser, &idEmail, &hashBytes, &expiresAt, &usedAt, &createdAt,
	); err != nil {
		return nil, err
	}
	hash, err := scanSHA256(hashBytes)
	if err != nil {
		return nil, err
	}
	return &mapper.EmailVerificationToken{
		ID:        id,
		IDUser:    idUser,
		IDEmail:   idEmail,
		Hash:      hash,
		ExpiresAt: expiresAt,
		UsedAt:    usedAt,
		CreatedAt: createdAt,
	}, nil
}

func (r EmailVerification) MarkUsed(ctx context.Context, id int64, usedAt time.Time) error {
	const query = `UPDATE email_verifications
	SET used_at = $2
	WHERE idemail_verifications = $1 AND used_at IS NULL;`
	_, err := r.db.Exec(ctx, query, id, usedAt)
	return err
}

func (r EmailVerification) RevokeAllByUser(ctx context.Context, uid int64, now time.Time) error {
	const query = `UPDATE email_verifications
	SET revoked_at = $1
	WHERE id_users = $2 AND revoked_at IS NULL;`
	_, err := r.db.Exec(ctx, query, now, uid)
	return err
}

func (r EmailVerification) GetActiveByUser(ctx context.Context, uid int64) (*mapper.EmailVerificationToken, error) {
	const query = `SELECT
		idemail_verifications, id_emails, token_hash, expires_at, used_at, created_at
	FROM email_verifications
	WHERE id_users = $1 AND revoked_at IS NULL AND expires_at > NOW()
	ORDER BY created_at DESC LIMIT 1;`
	var (
		id        int64
		idEmail   int64
		hashBytes []byte
		expiresAt time.Time
		usedAt    *time.Time
		createdAt time.Time
	)
	if err := r.db.QueryRow(ctx, query, uid).Scan(
		&id, &idEmail, &hashBytes, &expiresAt, &usedAt, &createdAt,
	); err != nil {
		return nil, err
	}
	hash, err := scanSHA256(hashBytes)
	if err != nil {
		return nil, err
	}
	return &mapper.EmailVerificationToken{
		ID:        id,
		IDUser:    uid,
		IDEmail:   idEmail,
		Hash:      hash,
		ExpiresAt: expiresAt,
		UsedAt:    usedAt,
		CreatedAt: createdAt,
	}, nil
}
