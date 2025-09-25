package persistence

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/test/utils"
)

func TestAuthenticationRepository(t *testing.T) {
	ctx := context.Background()
	cpf := utils.NewRandomString(10)
	password := utils.NewRandomString(10)
	sut := NewAuthenticateSpy()
	sut.SignIn(ctx, cpf, password)

	t.Run("should call SignIn once", func(t *testing.T) {
		if sut.CallCount != 1 {
			t.Errorf("expected call count to be 1, got %d", sut.CallCount)
		}
	})

	t.Run("should pass correct CPF and password", func(t *testing.T) {
		if sut.LastCPF != cpf && sut.LastPassword != password {
			t.Errorf("expected cpf and password to match, got %s and %s", sut.LastCPF, sut.LastPassword)
		}
	})

	t.Run("should return user not found error", func(t *testing.T) {
		sut.ShouldReturnUserNotFound = true
		_, err := sut.SignIn(ctx, cpf, password)
		if err == nil {
			t.Error("expected error to be not nil")
		}
	})

	t.Run("should return database error", func(t *testing.T) {
		sut.ShouldReturnUserNotFound = false
		sut.ShouldReturnError = true
		_, err := sut.SignIn(ctx, cpf, password)
		if err == nil {
			t.Error("expected error to be not nil")
		}
	})

	t.Run("should return authenticated user", func(t *testing.T) {
		sut.ShouldReturnError = false
		output, err := sut.SignIn(ctx, cpf, password)
		if err != nil {
			t.Error("expected error to be nil, got", err)
		}
		if output.ID != sut.CustomOutput.ID {
			t.Errorf("expected ID to match, got %s", output.ID)
		}
		if output.FirstName != sut.CustomOutput.FirstName {
			t.Errorf("expected FirstName to match, got %s", output.FirstName)
		}
		if output.LastName != sut.CustomOutput.LastName {
			t.Errorf("expected LastName to match, got %s", output.LastName)
		}
	})
}
