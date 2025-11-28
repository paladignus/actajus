package repository

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/stretchr/testify/assert"
)

type AuthenticationSpy struct{}

func (m *AuthenticationSpy) SignIn(ctx context.Context, email string) (dto.SignInOutput, error) {
	return dto.SignInOutput{}, nil
}

func (m *AuthenticationSpy) ValidatePassword(ctx context.Context, email, password string) error {
	return nil
}

func (m *AuthenticationSpy) FindEmailByCPF(ctx context.Context, cpf string) (dto.GetEmailByCPFOutput, error) {
	return dto.GetEmailByCPFOutput{}, nil
}

func (m *AuthenticationSpy) AccountIsActive(ctx context.Context, email string) (string, error) {
	return "active", nil
}

func (m *AuthenticationSpy) InvalidAllTokensByIDUser(ctx context.Context, id string) error {
	return nil
}

func (m *AuthenticationSpy) CreateRecoverPassword(ctx context.Context, email, id string) error {
	return nil
}

func TestAuthenticationInterface(t *testing.T) {
	ctx := context.Background()
	sut := &AuthenticationSpy{}
	signInOutput, err := sut.SignIn(ctx, "test@example.com")
	assert.NoError(t, err)
	err = sut.ValidatePassword(ctx, "test@example.com", "password")
	assert.NoError(t, err)
	emailOutput, err := sut.FindEmailByCPF(ctx, "12345678901")
	assert.NoError(t, err)
	_, err = sut.AccountIsActive(ctx, "test@example.com")
	assert.NoError(t, err)
	err = sut.InvalidAllTokensByIDUser(ctx, "user-id")
	assert.NoError(t, err)
	err = sut.CreateRecoverPassword(ctx, "test@example.com", "user-id")
	assert.NoError(t, err)
	assert.Empty(t, signInOutput)
	assert.Empty(t, emailOutput)
}

