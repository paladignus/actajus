package service

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestSignInInterface(t *testing.T) {
	ctx := context.Background()
	sut := &spy.SignIn{
		ExpectedOutput: dto.SignInOutput{
			IDUser:       123,
			FirstName:    "John",
			LastName:     "Doe",
			Email:        "john.doe@example.com",
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
		},
		ExpectedError: nil,
	}
	input := dto.SignInInput{
		CPF:      "12345678901",
		Password: "password123",
	}
	t.Run("should return the same IDUser, FirstName, Email, AccessToken and RefreshToken", func(t *testing.T) {
		output, err := sut.Execute(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, 123, output.IDUser)
		assert.Equal(t, "John", output.FirstName)
		assert.Equal(t, "john.doe@example.com", output.Email)
	})

	t.Run("should return an empty IDUser, FirstName, Email, AccessToken and RefreshToken", func(t *testing.T) {
		sut := &spy.SignIn{
			ExpectedOutput: dto.SignInOutput{},
			ExpectedError:  errors.New("authentication failed"),
		}
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "authentication failed", err.Error())
	})
}
