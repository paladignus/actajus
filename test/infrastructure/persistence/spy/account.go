// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type AuthenticateSpy struct {
	callCount     int
	findResult    dto.AuthenticateOutput
	findError     error
	validateError error
	// LastCPF                  string
	// LastPassword             string
}

func (a AuthenticateSpy) FindUserAccountByCPF(ctx context.Context, cpf string) (authenticated dto.AuthenticateOutput, err error) {
	return a.findResult, a.findError
}

func (a AuthenticateSpy) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
	return a.validateError
}

type MockLogger struct{}

func (l *MockLogger) Info(ctx context.Context, msg string, keyvals ...any)  {}
func (l *MockLogger) Warn(ctx context.Context, msg string, keyvals ...any)  {}
func (l *MockLogger) Error(ctx context.Context, msg string, keyvals ...any) {}
func (l *MockLogger) Debug(ctx context.Context, msg string, keyvals ...any) {}
func (l *MockLogger) With(...any) repository.Logger                         { return l }
func (l *MockLogger) WithError(err error) repository.Logger                 { return l }

// MockToken implementa gateway.Token
type MockToken struct {
	pair dto.TokenPair
	err  error

	calledWithID string
}

func (m *MockToken) GenerateTokenPair(idUser string) (dto.TokenPair, error) {
	m.calledWithID = idUser
	return m.pair, m.err
}

func (m *MockToken) RefreshAccessToken(refreshToken string) (dto.TokenPair, error) {
	return m.pair, m.err
}

func (m *MockToken) ValidateAccessToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

func (m *MockToken) ValidateRefreshToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

// type AuthenticateSpy struct {
// 	ShouldReturnError        bool
// 	ShouldReturnUserNotFound bool
// 	CustomOutput             dto.AuthenticateOutput
// 	CallCount                int
// 	LastCPF                  string
// 	LastPassword             string
// }
//
// func NewAccountSpy() *AuthenticateSpy {
// 	return &AuthenticateSpy{
// 		CustomOutput: dto.AuthenticateOutput{
// 			AccessToken:  "access_token",
// 			RefreshToken: "refresh_token",
// 			IDUser:       "1",
// 			FirstName:    "John",
// 			LastName:     "Doe",
// 			Email:        "tes@email.com",
// 			Roles:        []string{"root"},
// 		},
// 	}
// }
//
// func (a AuthenticateSpy) FindUserAccountByCPF(ctx context.Context, cpf string) (authenticated dto.AuthenticateOutput, err error) {
// 	return authenticated, err
// }
//
// func (a AuthenticateSpy) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
// 	return err
// }

// func (a *AuthenticateSpy) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedResponse, error) {
// 	a.CallCount++
// 	a.LastCPF = cpf
// 	a.LastPassword = password
// 	if a.ShouldReturnUserNotFound {
// 		return dto.AuthenticatedResponse{}, domain.ErrUserNotFound
// 	}
// 	if a.ShouldReturnError {
// 		return dto.AuthenticatedResponse{}, errors.New("spy database error")
// 	}
// 	return a.CustomOutput, nil
// }
