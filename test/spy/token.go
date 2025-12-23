// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

// type IPasswordResetToken interface {
// 	Create(ctx context.Context, token entity.PasswordResetToken) error
// 	FindByToken(ctx context.Context, token string) (entity.PasswordResetToken, error)
// 	MarkAsUsed(ctx context.Context, token string) error
// 	InvalidateUserTokens(ctx context.Context, idUser int) error
// }

type TokenSpy struct {
	CreateErr     error
	FindErr       error
	FindResult    entity.PasswordResetToken
	MaskErr       error
	InvalidateErr error
}

func (t *TokenSpy) Create(ctx context.Context, token entity.PasswordResetToken) error {
	return t.CreateErr
}

func (t *TokenSpy) FindByToken(ctx context.Context, token string) (entity.PasswordResetToken, error) {
	return t.FindResult, t.FindErr
}

func (t *TokenSpy) MarkAsUsed(ctx context.Context, token string) error {
	return t.MaskErr
}

func (t *TokenSpy) InvalidateUserTokens(ctx context.Context, idUser int) error {
	return t.InvalidateErr
}

type TokenServiceSpy struct {
	GenereateErr error
}

func (t *TokenServiceSpy) Generate() (string, error) {
	return "token", t.GenereateErr
}
