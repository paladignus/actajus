package usecase

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/usecase"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/test/infrastructure/persistence/spy"
	"github.com/paladignus/actajus/test/utils"
)

func TestSignIn(t *testing.T) {
	ctx := context.Background()
	cpf := utils.NewRandomString(10)
	password := utils.NewRandomString(10)
	req := dto.AuthenticateInput{CPF: cpf, Password: password}
	repository := spy.NewAccountSpy()
	config := config.Load()
	token := adapter.NewJWTAdapter(config.JWT)
	logger := adapter.NewDefaultLogger()
	sut := usecase.
		NewAuthenticate(repository, logger, token)
	sut.Execute(ctx, req)

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
		_, err := sut.Execute(ctx, req)
		if err == nil {
			t.Error("Expected error to be not nil")
		}
	})

	t.Run("should that an database error is returned", func(t *testing.T) {
		repository.ShouldReturnError = true
		_, err := sut.Execute(ctx, req)
		if err == nil {
			t.Error("Expected error to be not nil")
		}
	})

	t.Run("should return authenticated user", func(t *testing.T) {
		repository.ShouldReturnError = false
		repository.ShouldReturnUserNotFound = false
		output, err := sut.Execute(ctx, req)
		if err != nil {
			t.Error("Expected error to be nil")
		}
		if output.IDUser != repository.CustomOutput.IDUser {
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
