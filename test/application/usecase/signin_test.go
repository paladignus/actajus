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

	t.Run("should corrects properties from repository", func(t *testing.T) {
		if repository.LastCPF != cpf || repository.LastPassword != password {
			t.Error("Expected cpf and password to match")
		}
	})

	t.Run("should call repository only once", func(t *testing.T) {
		if repository.CallCount != 1 {
			t.Error("Expected callsCount to be 1")
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
