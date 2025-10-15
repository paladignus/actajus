package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/usecase"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/test/infrastructure/persistence/spy"
)

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()
	validCPF := "11144477735"
	input := dto.AuthenticateInput{
		CPF:      "123",
		Password: "whatever",
	}
	account := &spy.AuthenticateSpy{}
	logger := &spy.MockLogger{}
	token := &spy.MockToken{}

	t.Run("should return error ErrInvalidCPF", func(t *testing.T) {
		sut := usecase.NewAuthenticate(account, logger, token)
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, exception.ErrInvalidCPF) {
			t.Errorf("expected to be '%v', but got '%v'", exception.ErrInvalidCPF, err)
		}
	})

	t.Run("should return error if password user not found", func(t *testing.T) {
		wantErr := errors.New("not found from repo")
		input.CPF = validCPF
		account.FindError = wantErr
		sut := usecase.NewAuthenticate(account, logger, token)
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, wantErr) {
			t.Errorf("expected to be '%v', but got '%v'", wantErr, err)
		}
	})

	t.Run("should return error ErrInvalidPassword", func(t *testing.T) {
		wantErr := errors.New("invalid credentials")
		account.FindResult.IDUser = "user-1"
		account.ValidateError = wantErr
		account.FindError = nil
		sut := usecase.NewAuthenticate(account, logger, token)
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, wantErr) {
			t.Errorf("expected to be '%v', but got '%v'", wantErr, err)
		}
	})
}
