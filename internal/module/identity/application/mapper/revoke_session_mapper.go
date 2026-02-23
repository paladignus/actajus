// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/model"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
)

func (m *AuthMapper) RevokeInputToNormalized(input dto.RevokeSessionCommand) (model.RevokeNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return model.RevokeNormalized{}, err
	}
	return model.RevokeNormalized{IDSession: input.IDSession}, nil
}
