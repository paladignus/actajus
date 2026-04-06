// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/command"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
)

func (m *AuthMapper) LogoutInputToNormalized(input command.LogoutCommand) (LogoutNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return LogoutNormalized{}, err
	}
	return LogoutNormalized{IDSession: input.IDSession}, nil
}

func (m *AuthMapper) LogoutAllInputToNormalized(input command.LogoutAllCommand) (LogoutAllNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return LogoutAllNormalized{}, err
	}
	return LogoutAllNormalized{IDUser: input.IDUser}, nil
}
