package usecase

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/application/usecase"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/spy"
)

func TestSignIn(t *testing.T) {
	ctx := context.Background()
	cpf := NewRandomString(10)
	password := NewRandomString(10)
	repository := spy.NewAuthenticateSpy()
	sut := usecase.NewSignIn(repository)
	sut.Execute(ctx, cpf, password)

	t.Run("should corrects params from repository", func(t *testing.T) {
		if repository.LastCPF != cpf || repository.LastPassword != password {
			t.Error("Expected cpf and password to match")
		}
	})

	t.Run("should call repository only once", func(t *testing.T) {
		if repository.CallCount != 1 {
			t.Error("Expected callsCount to be 1")
		}
	})

	t.Run("should return user not found error", func(t *testing.T) {
		repository.ShouldReturnUserNotFound = true
		_, err := sut.Execute(ctx, cpf, password)
		if err == nil {
			t.Error("Expected error to be not nil")
		}
	})

	t.Run("should that an database error is returned", func(t *testing.T) {
		repository.ShouldReturnError = true
		_, err := sut.Execute(ctx, cpf, password)
		if err == nil {
			t.Error("Expected error to be not nil")
		}
	})

	t.Run("should return authenticated user", func(t *testing.T) {
		repository.ShouldReturnError = false
		repository.ShouldReturnUserNotFound = false
		output, err := sut.Execute(ctx, cpf, password)
		if err != nil {
			t.Error("Expected error to be nil")
		}
		if output.ID != repository.CustomOutput.ID {
			t.Error("Expected ID to match")
		}
		if output.FirstName != repository.CustomOutput.FirstName {
			t.Error("Expected FirstName to match")
		}
		if output.LastName != repository.CustomOutput.LastName {
			t.Error("Expected LastName to match")
		}
	})
}

func NewRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
