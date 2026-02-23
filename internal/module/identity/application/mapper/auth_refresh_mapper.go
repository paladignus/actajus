// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/model"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
)

func (m *AuthMapper) RefreshInputToNormalized(input dto.RefreshCommand) (model.RefreshNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return model.RefreshNormalized{}, err
	}
	return model.RefreshNormalized{
		IDSession:    input.IDSession,
		RefreshToken: input.RefreshToken,
		IP:           input.IP,
		UserAgent:    input.UserAgent,
	}, nil
}
