// Package postgres
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type PasswordResetToken struct {
	db *database.DB
}

func NewPasswordResetToken(db *database.DB) PasswordResetToken {
	return PasswordResetToken{db}
}

func (p PasswordResetToken) Create(ctx context.Context, token entity.PasswordResetToken) error {
	sql := `INSERT INTO password_reset (id_users, token, expires_at) VALUES ($1, $2, $3);`
	_, err := p.db.Pool.Exec(ctx, sql, token.IDUser, token.Token, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("database error while saving password reset token for user ID %d: %w", token.IDUser, err)
	}
	return nil
}

func (p PasswordResetToken) FindByToken(ctx context.Context, token string) (pr entity.PasswordResetToken, err error) {
	sql := `
		SELECT
			idpassword_reset,
			id_users,
			token,
			expires_at,
			used_at
		FROM password_reset pr
		JOIN documents d ON pr.id_users = d.id_people
		WHERE token = $1 AND used_at IS NULL;`
	if err = p.db.Pool.QueryRow(ctx, sql, token).Scan(
		&pr.IDPasswordReset,
		&pr.IDUser,
		&pr.Token,
		&pr.ExpiresAt,
		&pr.UsedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr, fmt.Errorf("invalid token %s: %w", token, exception.ErrInvalidToken)
		}
		return pr, fmt.Errorf("database error while finding token %s: %w", token, err)
	}
	return pr, nil
}

func (p PasswordResetToken) MarkAsUsed(ctx context.Context, token string) error {
	sql := `UPDATE password_reset SET used_at = now() WHERE token = $1 AND used_at IS NULL;`
	_, err := p.db.Pool.Exec(ctx, sql, token)
	if err != nil {
		return fmt.Errorf("database error while marking token %s as used: %w", token, err)
	}
	return nil
}

func (p PasswordResetToken) InvalidateUserTokens(ctx context.Context, idUser int) error {
	sql := `UPDATE password_reset SET used_at = now() WHERE id_users = $1 AND used_at IS NULL;`
	_, err := p.db.Pool.Exec(ctx, sql, idUser)
	if err != nil {
		return fmt.Errorf("database error while invalidating tokens for user ID %d: %w", idUser, err)
	}
	return nil
}
