// Package repository
package repository

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

type PasswordResetTokenSpy struct {
	PasswordResetToken string
}

func (p *PasswordResetTokenSpy) Create(ctx context.Context, token entity.PasswordResetToken) error {
	return nil
}

func (p *PasswordResetTokenSpy) FindByToken(ctx context.Context, token string) (entity.PasswordResetToken, error) {
	return entity.PasswordResetToken{}, nil
}

func (p *PasswordResetTokenSpy) MarkAsUsed(ctx context.Context, token string) error {
	return nil
}

func (p *PasswordResetTokenSpy) InvalidateUserTokens(ctx context.Context, idUser int) error {
	return nil
}

func TestPasswordResetTokenInterface(t *testing.T) {
	ctx := context.Background()
	sut := &PasswordResetTokenSpy{}
	err := sut.Create(ctx, entity.PasswordResetToken{})
	assert.NoError(t, err)
	passwordResetToken, err := sut.FindByToken(ctx, "token")
	assert.NoError(t, err)
	assert.Equal(t, passwordResetToken, entity.PasswordResetToken{})
	err = sut.MarkAsUsed(ctx, "token")
	assert.NoError(t, err)
	err = sut.InvalidateUserTokens(ctx, 1)
	assert.NoError(t, err)
}
