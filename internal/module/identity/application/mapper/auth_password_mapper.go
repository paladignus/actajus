// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
)

func (m *AuthMapper) ChangePasswordInputToNormalized(input dto.ChangePasswordCommand) (ChangePasswordNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	// if input.UserID <= 0 && !hasViolation(vs, "id_user") {
	// 	vs = append(vs, validation.Violation{Path: "id_user", Code: "invalid", Meta: map[string]string{}})
	// }
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return ChangePasswordNormalized{}, err
	}
	return ChangePasswordNormalized{
		IDUser:          input.IDUser,
		CurrentPassword: input.CurrentPassword,
		NewPassword:     input.NewPassword,
	}, nil
}

func (m *AuthMapper) RequestPasswordResetInputToNormalized(input dto.RequestPasswordResetCommand) (RequestPasswordResetNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	// email := vo.Email(input.Email)
	// if !email.IsEmpty() && !email.IsValid() && !hasViolation(vs, "email") {
	// 	vs = append(vs, validation.Violation{Path: "email", Code: "invalid", Meta: map[string]string{"format": "email"}})
	// }
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return RequestPasswordResetNormalized{}, err
	}
	return RequestPasswordResetNormalized{Email: input.Email}, nil
}

func (m *AuthMapper) ConfirmPasswordResetInputToNormalized(input dto.ConfirmPasswordResetCommand) (ConfirmPasswordResetNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	// if input.IDReset <= 0 && !hasViolation(vs, "id_reset") {
	// 	vs = append(vs, validation.Violation{Path: "id_reset", Code: "invalid", Meta: map[string]string{}})
	// }
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return ConfirmPasswordResetNormalized{}, err
	}
	return ConfirmPasswordResetNormalized{
		IDReset:     input.IDReset,
		ResetToken:  input.ResetToken,
		NewPassword: input.NewPassword,
	}, nil
}
