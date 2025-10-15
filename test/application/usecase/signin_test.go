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
	// validCPF := "11144477735"

	t.Run("should return error ErrInvalidCPF", func(t *testing.T) {
		input := dto.AuthenticateInput{
			CPF:      "123",
			Password: "whatever",
		}
		account := &spy.AuthenticateSpy{}
		logger := &spy.MockLogger{}
		token := &spy.MockToken{}

		auth := usecase.NewAuthenticate(account, logger, token)
		_, err := auth.Execute(ctx, input)
		if !errors.Is(err, exception.ErrInvalidCPF) {
			t.Fatalf("esperava ErrInvalidCPF, recebeu: %v", err)
		}
	})
}

// func TestAuthenticate(t *testing.T) {
// 	ctx := context.Background()
// 	cpf := utils.NewRandomString(10)
// 	password := utils.NewRandomString(10)
// 	repository := spy.NewAccountSpy()
// 	config := config.Load()
// 	token := adapter.NewJWTAdapter(config.JWT)
// 	logger := adapter.NewDefaultLogger()
// 	sut := usecase.NewAuthenticate(repository, logger, token)
// 	sut.Execute(ctx, dto.AuthenticateInput{CPF: cpf, Password: password})
//
// 	// t.Run("should corrects params from repository", func(t *testing.T) {
// 	// 	t.SkipNow()
// 	// 	if repository.LastCPF != cpf || repository.LastPassword != password {
// 	// 		t.Error("Expected cpf and password to match")
// 	// 	}
// 	// })
//
// 	t.Run("should call repository only once", func(t *testing.T) {
// 		if repository.CallCount != 1 {
// 			t.Error("Expected callsCount to be 1")
// 		}
// 	})
//
// 	// t.Run("should return user not found error", func(t *testing.T) {
// 	// 	t.SkipNow()
// 	// 	repository.ShouldReturnUserNotFound = true
// 	// 	_, err := sut.Execute(ctx, req)
// 	// 	if err == nil {
// 	// 		t.Error("Expected error to be not nil")
// 	// 	}
// 	// })
//
// 	// t.Run("should that an database error is returned", func(t *testing.T) {
// 	// 	t.SkipNow()
// 	// 	repository.ShouldReturnError = true
// 	// 	_, err := sut.Execute(ctx, req)
// 	// 	if err == nil {
// 	// 		t.Error("Expected error to be not nil")
// 	// 	}
// 	// })
//
// 	// t.Run("should return authenticated user", func(t *testing.T) {
// 	// 	t.SkipNow()
// 	// 	repository.ShouldReturnError = false
// 	// 	repository.ShouldReturnUserNotFound = false
// 	// 	output, err := sut.Execute(ctx, req)
// 	// 	if err != nil {
// 	// 		t.Error("Expected error to be nil")
// 	// 	}
// 	// 	if output.IDUser != repository.CustomOutput.IDUser {
// 	// 		t.Error("Expected ID to match")
// 	// 	}
// 	// 	if output.FirstName != repository.CustomOutput.FirstName {
// 	// 		t.Error("Expected FirstName to match")
// 	// 	}
// 	// 	if output.LastName != repository.CustomOutput.LastName {
// 	// 		t.Error("Expected LastName to match")
// 	// 	}
// 	// })
// }
