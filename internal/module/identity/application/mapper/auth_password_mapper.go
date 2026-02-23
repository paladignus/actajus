// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/model"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

func (m *AuthMapper) ChangePasswordInputToNormalized(input dto.ChangePasswordCommand) (model.ChangePasswordNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	// if input.UserID <= 0 && !hasViolation(vs, "id_user") {
	// 	vs = append(vs, validation.Violation{Path: "id_user", Code: "invalid", Meta: map[string]string{}})
	// }
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return model.ChangePasswordNormalized{}, err
	}
	return model.ChangePasswordNormalized{
		IDUser:          input.IDUser,
		CurrentPassword: input.CurrentPassword,
		NewPassword:     input.NewPassword,
	}, nil
}

func (m *AuthMapper) RequestPasswordResetInputToNormalized(input dto.RequestPasswordResetCommand) (model.RequestPasswordResetNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	email := vo.Email(input.Email)
	// if !email.IsEmpty() && !email.IsValid() && !hasViolation(vs, "email") {
	// 	vs = append(vs, validation.Violation{Path: "email", Code: "invalid", Meta: map[string]string{"format": "email"}})
	// }
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return model.RequestPasswordResetNormalized{}, err
	}
	return model.RequestPasswordResetNormalized{Email: email}, nil
}

func (m *AuthMapper) ConfirmPasswordResetInputToNormalized(input dto.ConfirmPasswordResetCommand) (model.ConfirmPasswordResetNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	// if input.IDReset <= 0 && !hasViolation(vs, "id_reset") {
	// 	vs = append(vs, validation.Violation{Path: "id_reset", Code: "invalid", Meta: map[string]string{}})
	// }
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return model.ConfirmPasswordResetNormalized{}, err
	}
	return model.ConfirmPasswordResetNormalized{
		IDReset:     input.IDReset,
		ResetToken:  input.ResetToken,
		NewPassword: input.NewPassword,
	}, nil
}
