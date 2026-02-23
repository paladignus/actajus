// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/model"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
)

func (m *AuthMapper) LogoutInputToNormalized(input dto.LogoutCommand) (model.LogoutNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return model.LogoutNormalized{}, err
	}
	return model.LogoutNormalized{IDSession: input.IDSession}, nil
}

func (m *AuthMapper) LogoutAllInputToNormalized(input dto.LogoutAllCommand) (model.LogoutAllNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return model.LogoutAllNormalized{}, err
	}
	return model.LogoutAllNormalized{IDUser: input.IDUser}, nil
}
