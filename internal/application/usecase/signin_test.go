package usecase

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	ctx := context.Background()
	validCPF := "11144477735"
	input := dto.SignInInput{
		CPF:      "123",
		Password: "whatever",
	}
	user := spy.NewUser()
	logger := &spy.Logger{}
	token := &spy.JWT{}

	t.Run("should return error ErrInvalidCPF", func(t *testing.T) {
		sut := NewSignIn(user, token)
		logger.On("Warn", ctx, "invalid cpf format provided", "cpf", input.CPF).Once()
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid CPF format 123 in sign in use case")
	})

	t.Run("should return error if password user not found", func(t *testing.T) {
		wantErr := errors.New("not found from repo")
		input.CPF = validCPF
		user.FindError = wantErr
		sut := NewSignIn(user, token)
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sign in use case failed for CPF 11144477735")
	})

	t.Run("should return error if fails generating tokens", func(t *testing.T) {
		wantErr := adapter.ErrBuildToken
		input.Password = "!M@r1L0$n4"
		user.FindError = nil
		token.Err = wantErr
		sut := NewSignIn(user, token)
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sign in use case failed to generate token pair for user ID")
	})

	t.Run("should return all correct values", func(t *testing.T) {
		token.Pair.AccessToken = "access-token-xyz"
		token.Pair.RefreshToken = "refresh-token-abc"
		token.Err = nil
		sut := NewSignIn(user, token)
		output, err := sut.Execute(ctx, input)
		if err != nil {
			t.Fatalf("should not return error, got: %v", err)
		}
		if output.IDUser != user.FindResult.UserDTO.IDUser {
			t.Errorf("expected IDUser to be '%d', but got '%d'", user.FindResult.UserDTO.IDUser, output.IDUser)
		}
		if output.AccessToken != token.Pair.AccessToken {
			t.Errorf("expected AccessToken to be '%s', but got '%s'", token.Pair.AccessToken, output.AccessToken)
		}
		if output.RefreshToken != token.Pair.RefreshToken {
			t.Errorf("expected RefreshToken to be '%s', but got '%s'", token.Pair.RefreshToken, output.RefreshToken)
		}
		idUser := strconv.Itoa(user.FindResult.UserDTO.IDUser)
		if token.CalledWithID != idUser {
			t.Errorf("expected CalledWithID to be '%d', but got '%s'", user.FindResult.UserDTO.IDUser, token.CalledWithID)
		}
	})
}
