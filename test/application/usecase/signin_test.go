package usecase_test

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
	repository := spy.AuthenticateSpy{}
	sut := usecase.NewSignIn(&repository)
	sut.Execute(ctx, cpf, password)

	t.Run("should corrects properties from repository", func(t *testing.T) {
		if repository.CPF != cpf || repository.Password != password {
			t.Error("Expected cpf and password to match")
		}
	})

	t.Run("should call repository only once", func(t *testing.T) {
		if repository.CallsCount != 1 {
			t.Error("Expected callsCount to be 1")
		}
	})

	t.Run("should return user id", func(t *testing.T) {
		people, _ := sut.Execute(ctx, cpf, password)
		if people.ID == "" {
			t.Error("Expected userID to be not empty")
		}
	})

	t.Run("should return an error if cpf id empty", func(t *testing.T) {
		_, err := sut.Execute(ctx, "", password)
		if err == nil {
			t.Error("Expected error to be nil")
		}
	})

	t.Run("should return an error if password id empty", func(t *testing.T) {
		_, err := sut.Execute(ctx, cpf, "")
		if err == nil {
			t.Error("Expected error to be nil")
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
