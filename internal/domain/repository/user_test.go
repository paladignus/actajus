package repository

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/stretchr/testify/assert"
)

type UserSpy struct{}

func (u *UserSpy) AuthenticationByCPF(ctx context.Context, input command.SignInCommand) (user readmodel.SignInReadModel, err error) {
	return readmodel.SignInReadModel{}, nil
}

func (u *UserSpy) ValidatePassword(ctx context.Context, email, password string) error {
	return nil
}

func (u *UserSpy) FindEmailByCPF(ctx context.Context, cpf string) (readmodel.GetEmailByCPFReadModel, error) {
	return readmodel.GetEmailByCPFReadModel{}, nil
}

func (u *UserSpy) FindIDUserByEmail(ctx context.Context, email string) (int, error) {
	return 0, nil
}

func (u *UserSpy) UpdatePassword(ctx context.Context, idUser int, password, cpf string) error {
	return nil
}

func TestUserInterface(t *testing.T) {
	ctx := context.Background()
	sut := &UserSpy{}
	signInOutput, err := sut.AuthenticationByCPF(ctx, command.SignInCommand{})
	assert.NoError(t, err)
	assert.Empty(t, signInOutput)
	err = sut.ValidatePassword(ctx, "test@example.com", "password")
	assert.NoError(t, err)
	emailOutput, err := sut.FindEmailByCPF(ctx, "12345678901")
	assert.NoError(t, err)
	assert.Empty(t, emailOutput)
	idUser, err := sut.FindIDUserByEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, idUser, 0)
	err = sut.UpdatePassword(ctx, 1, "password", "11144477735")
	assert.NoError(t, err)
}
