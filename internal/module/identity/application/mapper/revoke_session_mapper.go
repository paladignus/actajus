// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/command"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
)

func (m *AuthMapper) RevokeInputToNormalized(input command.RevokeSessionCommand) (RevokeNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return RevokeNormalized{}, err
	}
	return RevokeNormalized{IDSession: input.IDSession}, nil
}
