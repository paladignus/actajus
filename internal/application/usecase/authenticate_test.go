package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/test/spy"
)

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()
	validCPF := "11144477735"
	input := dto.AuthenticateInput{
		CPF:      "123",
		Password: "whatever",
	}
	account := spy.NewAuthenticateSpy()
	logger := &spy.MockLogger{}
	token := &spy.MockToken{}

	t.Run("should return error ErrInvalidCPF", func(t *testing.T) {
		sut := NewAuthenticate(account, logger, token)
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, exception.ErrInvalidCPF) {
			t.Errorf("expected to be '%v', but got '%v'", exception.ErrInvalidCPF, err)
		}
	})

	t.Run("should return error if password user not found", func(t *testing.T) {
		wantErr := errors.New("not found from repo")
		input.CPF = validCPF
		account.FindError = wantErr
		sut := NewAuthenticate(account, logger, token)
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
		sut := NewAuthenticate(account, logger, token)
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, wantErr) {
			t.Errorf("expected to be '%v', but got '%v'", wantErr, err)
		}
	})

	t.Run("should return error if fails generating tokens", func(t *testing.T) {
		wantErr := errors.New("token generation failed")
		input.Password = "!M@r1L0$n4"
		account.FindResult.IDUser = "user-123"
		account.ValidateError = nil
		token.Err = wantErr
		sut := NewAuthenticate(account, logger, token)
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, wantErr) {
			t.Errorf("expected to be '%v', but got '%v'", wantErr, err)
		}
	})

	t.Run("should return all correct values", func(t *testing.T) {
		token.Pair.AccessToken = "access-token-xyz"
		token.Pair.RefreshToken = "refresh-token-abc"
		token.Err = nil
		sut := NewAuthenticate(account, logger, token)
		output, err := sut.Execute(ctx, input)
		if err != nil {
			t.Fatalf("should not return error, got: %v", err)
		}
		if output.IDUser != account.FindResult.IDUser {
			t.Errorf("expected IDUser to be '%s', but got '%s'", account.FindResult.IDUser, output.IDUser)
		}
		if output.AccessToken != token.Pair.AccessToken {
			t.Errorf("expected AccessToken to be '%s', but got '%s'", token.Pair.AccessToken, output.AccessToken)
		}
		if output.RefreshToken != token.Pair.RefreshToken {
			t.Errorf("expected RefreshToken to be '%s', but got '%s'", token.Pair.RefreshToken, output.RefreshToken)
		}
		if token.CalledWithID != account.FindResult.IDUser {
			t.Errorf("expected CalledWithID to be '%s', but got '%s'", account.FindResult.IDUser, token.CalledWithID)
		}
	})
}
