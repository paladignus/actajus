// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
)

type AuthMapper struct {
	validator *validation.Validator
}

func NewAuthMapper(v *validation.Validator) *AuthMapper {
	return &AuthMapper{v}
}

func (m *AuthMapper) LoginInputToNormalized(input dto.LoginCommand) (LoginNormalized, error) {
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
