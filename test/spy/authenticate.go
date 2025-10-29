// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type AuthenticateSpy struct {
	CallCount     int
	FindResult    dto.AuthenticateOutput
	FindError     error
	ValidateError error
}

func NewAuthenticateSpy() *AuthenticateSpy {
	return &AuthenticateSpy{}
}

func (a *AuthenticateSpy) FindUserAccountByCPF(ctx context.Context, cpf string) (authenticated dto.AuthenticateOutput, err error) {
	return a.FindResult, a.FindError
}

func (a *AuthenticateSpy) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
	return a.ValidateError
}

func (a *AuthenticateSpy) GetEmailByCPF(ctx context.Context, cpf string) (string, error) {
	return a.FindResult.Email, a.FindError
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
	Pair       dto.TokenPair
	TokenReset dto.TokenRecoverPassword
	Err        error

	CalledWithID string
}

func (m *MockToken) GenerateTokenPair(idUser string) (dto.TokenPair, error) {
	m.CalledWithID = idUser
	return m.Pair, m.Err
}

func (m *MockToken) RefreshAccessToken(refreshToken string) (dto.TokenPair, error) {
	return m.Pair, m.Err
}

func (m *MockToken) ValidateAccessToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

func (m *MockToken) ValidateRefreshToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

func (m *MockToken) GenerateRecoverPasswordToken(idUser string) (dto.TokenRecoverPassword, error) {
	return m.TokenReset, m.Err
}

func (m *MockToken) ValidateRecoverPasswordToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}
