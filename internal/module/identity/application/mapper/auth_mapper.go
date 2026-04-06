// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/command"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
)

type AuthMapper struct {
	validator *validation.Validator
}

func NewAuthMapper(v *validation.Validator) *AuthMapper {
	return &AuthMapper{v}
}

func (m *AuthMapper) LoginInputToNormalized(input command.LoginCommand) (LoginNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return LoginNormalized{}, err
	}
	return LoginNormalized{
		Email:     input.Email,
		Password:  input.Password,
		IP:        input.IP,
		UserAgent: input.UserAgent,
	}, nil
}

func (m *AuthMapper) RegisterInputToNormalized(input command.RegisterCommand) (RegisterNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return RegisterNormalized{}, err
	}
	return RegisterNormalized{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Birthday:  input.Birthday,
		GenderID:  input.GenderID,
		Email:     input.Email,
		Password:  input.Password,
	}, nil
}

func (m *AuthMapper) ConfirmEmailVerificationInputToNormalized(input command.ConfirmEmailVerificationCommand) (ConfirmEmailVerificationNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return ConfirmEmailVerificationNormalized{}, err
	}
	return ConfirmEmailVerificationNormalized{
		IDVerification: input.IDVerification,
		Token:          input.Token,
	}, nil
}
