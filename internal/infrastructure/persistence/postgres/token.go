// Package postgres
package postgres

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type Token struct {
	db *database.DB
}

func NewToken(db *database.DB) Token {
	return Token{db}
}

func (t Token) Create(ctx context.Context, token entity.PasswordResetToken) (err error) {
	sql := `INSERT INTO password_reset (id_people, token, expires_at) VALUES ($1, $2, $3);`
	_, err = t.db.Pool.Exec(ctx, sql, token.IDAccount, token.Token, token.ExpiresAt)
	return err
}

func (t Token) FindByToken(ctx context.Context, token string) (entity.PasswordResetToken, error) {
	return entity.PasswordResetToken{}, nil
}

func (t Token) MarkAsUsed(ctx context.Context, token string) error {
	return nil
}

func (t Token) InvalidateUserTokens(ctx context.Context, idUser int) error {
	return nil
}
