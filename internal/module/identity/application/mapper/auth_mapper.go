// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/model"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
)

type AuthMapper struct {
	validator *validation.Validator
}

func NewAuthMapper(v *validation.Validator) *AuthMapper {
	return &AuthMapper{v}
}

func (m *AuthMapper) LoginInputToNormalized(input dto.LoginCommand) (model.LoginNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	email := vo.Email(input.Email)
	if shouldAddEmailViolation(vs) {
		if !email.IsEmpty() && !email.IsValid() {
			vs = append(vs, validation.Violation{
				Path: "email",
				Code: "invalid",
				Meta: map[string]string{"format": "email"},
			})
		}
	}
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return model.LoginNormalized{}, err
	}
	return model.LoginNormalized{
		Email:     email,
		Password:  input.Password,
		IP:        input.IP,
		UserAgent: input.UserAgent,
	}, nil
}
