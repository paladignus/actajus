// Package spy
package spy

import (
	"context"
	"errors"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/domainerrors"
)

// type AuthenticateSpy struct {
// 	CPF        string
// 	Password   string
// 	CallsCount int
// }
//
// func (a *AuthenticateSpy) SignIn(ctx context.Context, cpf, password string) (dto.AuthenticatedOutput, error) {
// 	authenticated := dto.AuthenticatedOutput{}
// 	if cpf == "" || password == "" {
// 		return authenticated, errors.New("cpf and password cannot be empty")
// 	}
// 	a.CPF = cpf
// 	a.Password = password
// 	a.CallsCount++
// 	authenticated.ID = "1"
// 	authenticated.FirstName = "John"
// 	authenticated.LastName = "Doe"
// 	return authenticated, nil
// }

// AuthenticateSpy é uma implementação mock para testes unitários
type AuthenticateSpy struct {
	// Configuração do comportamento do mock
	ShouldReturnError        bool
	ShouldReturnUserNotFound bool
	CustomOutput             dto.AuthenticatedOutput
	CallCount                int
	LastCPF                  string
	LastPassword             string
}

// NewAuthenticateSpy cria uma nova instância do mock
func NewAuthenticateSpy() *AuthenticateSpy {
	return &AuthenticateSpy{
		CustomOutput: dto.AuthenticatedOutput{
			ID:        "1",
			FirstName: "Mock",
			LastName:  "User",
		},
	}
}

// SignIn implementa a interface do Authenticate para testes
func (a *AuthenticateSpy) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
	a.CallCount++
	a.LastCPF = cpf
	a.LastPassword = password

	// Simula erro de usuário não encontrado
	if a.ShouldReturnUserNotFound {
		return dto.AuthenticatedOutput{}, domainerrors.ErrUserNotFound
	}

	// Simula outros erros
	if a.ShouldReturnError {
		return dto.AuthenticatedOutput{}, errors.New("mock database error")
	}

	// Retorna o output customizado ou padrão
	return a.CustomOutput, nil
}

// AuthenticateWithFixedResponseSpy cria um mock com resposta fixa
type AuthenticateWithFixedResponseSpy struct {
	FixedResponse dto.AuthenticatedOutput
	FixedError    error
	CallCount     int
}

func NewAuthenticateWithFixedResponseSpy(response dto.AuthenticatedOutput, err error) *AuthenticateWithFixedResponseSpy {
	return &AuthenticateWithFixedResponseSpy{
		FixedResponse: response,
		FixedError:    err,
	}
}

func (m *AuthenticateWithFixedResponseSpy) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
	m.CallCount++
	return m.FixedResponse, m.FixedError
}
