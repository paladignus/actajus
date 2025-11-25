package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/test/spy"
)

func TestSignIn(t *testing.T) {
	ctx := context.Background()
	validCPF := "11144477735"
	input := dto.SignInInput{
		CPF:      "123",
		Password: "whatever",
	}
	authentication := spy.NewAuthenticationSpy()
	logger := &spy.SpyLogger{}
	token := &spy.SpyToken{}

	t.Run("should return error ErrInvalidCPF", func(t *testing.T) {
		sut := NewSignIn(authentication, logger, token)
		logger.On("Warn", ctx, "invalid cpf format provided", "cpf", input.CPF).Once()
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, exception.ErrInvalidCPF) {
			t.Errorf("expected to be '%v', but got '%v'", exception.ErrInvalidCPF, err)
		}
	})

	t.Run("should return error if password user not found", func(t *testing.T) {
		wantErr := errors.New("not found from repo")
		input.CPF = validCPF
		authentication.FindError = wantErr
		sut := NewSignIn(authentication, logger, token)
		logger.On("Warn", ctx, "user not found during authentication",
			"cpf", input.CPF,
			"error", wantErr).Once()
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, wantErr) {
			t.Errorf("expected to be '%v', but got '%v'", wantErr, err)
		}
	})

	t.Run("should return error ErrInvalidPassword", func(t *testing.T) {
		wantErr := errors.New("invalid credentials")
		authentication.ValidateError = wantErr
		authentication.FindError = nil
		sut := NewSignIn(authentication, logger, token)
		logger.On("Warn", ctx, "invalid credentials provided",
			"cpf", input.CPF,
			"id_person", authentication.FindResult.Authentication.IDUser).Once()
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, wantErr) {
			t.Errorf("expected to be '%v', but got '%v'", wantErr, err)
		}
	})

	t.Run("should return error if fails generating tokens", func(t *testing.T) {
		wantErr := adapter.ErrBuildToken
		input.Password = "!M@r1L0$n4"
		authentication.ValidateError = nil
		token.Err = wantErr
		sut := NewSignIn(authentication, logger, token)
		logger.On("Info", ctx, "person authenticated successfully",
			"cpf", input.CPF,
			"id_person", authentication.FindResult.Authentication.IDUser).Once()
		logger.On("Error", ctx, "failed to generate token pair", "error", wantErr).Once()
		_, err := sut.Execute(ctx, input)
		if !errors.Is(err, wantErr) {
			t.Errorf("expected to be '%v', but got '%v'", wantErr, err)
		}
	})

	t.Run("should return all correct values", func(t *testing.T) {
		token.Pair.AccessToken = "access-token-xyz"
		token.Pair.RefreshToken = "refresh-token-abc"
		token.Err = nil
		sut := NewSignIn(authentication, logger, token)
		logger.On("Info", ctx, "person authenticated successfully",
			"cpf", input.CPF,
			"id_person", authentication.FindResult.Authentication.IDUser).Once()
		output, err := sut.Execute(ctx, input)
		if err != nil {
			t.Fatalf("should not return error, got: %v", err)
		}
		if output.IDUser != authentication.FindResult.Authentication.IDUser {
			t.Errorf("expected IDUser to be '%s', but got '%s'", authentication.FindResult.Authentication.IDUser, output.IDUser)
		}
		if output.AccessToken != token.Pair.AccessToken {
			t.Errorf("expected AccessToken to be '%s', but got '%s'", token.Pair.AccessToken, output.AccessToken)
		}
		if output.RefreshToken != token.Pair.RefreshToken {
			t.Errorf("expected RefreshToken to be '%s', but got '%s'", token.Pair.RefreshToken, output.RefreshToken)
		}
		if token.CalledWithID != authentication.FindResult.Authentication.IDUser {
			t.Errorf("expected CalledWithID to be '%s', but got '%s'", authentication.FindResult.Authentication.IDUser, token.CalledWithID)
		}
	})
}
