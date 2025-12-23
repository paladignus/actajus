// Package usecase
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestRecoverPassword(t *testing.T) {
	ctx := context.Background()
	user := spy.NewUser()
	token := &spy.TokenSpy{}
	// service := security.NewCryptoTokenGenerator()
	service := &spy.TokenServiceSpy{}
	publisher := &spy.Publisher{}
	sut := NewRequestPasswordReset(
		user, token, service, publisher)
	input := dto.RequestPasswordResetInput{Email: "email"}

	t.Run("should return error if email is invalid", func(t *testing.T) {
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid email format email in request password reset use case")
	})

	t.Run("should return an error if it fails to find the id user", func(t *testing.T) {
		wantErr := errors.New("failed to find id user")
		user.FindError = wantErr
		input.Email = "email@example.com.br"
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "request password reset use case failed to find user by email email@example.com.br")
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("should return an error if it fails to invalidate previous token", func(t *testing.T) {
		wantErr := errors.New("failed to invalidate previous token")
		token.InvalidateErr = wantErr
		user.FindError = nil
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
	})

	t.Run("should return an errors if generate token its failed", func(t *testing.T) {
		wantErr := adapter.ErrBuildToken
		service.GenereateErr = wantErr
		token.InvalidateErr = nil
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("should return an error if it fails to save token", func(t *testing.T) {
		wantErr := errors.New("failed to save token")
		service.GenereateErr = nil
		token.CreateErr = wantErr
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("should return an error if it fails publisher event", func(t *testing.T) {
		wantErr := errors.New("failed to publisher event")
		publisher.Err = wantErr
		token.CreateErr = nil
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("should successful to recover password", func(t *testing.T) {
		publisher.Err = nil
		err := sut.Execute(ctx, input)
		assert.NoError(t, err)
	})
}
