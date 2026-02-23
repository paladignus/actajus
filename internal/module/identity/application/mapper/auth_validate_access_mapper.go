// Package mapper
package mapper

import (
	"strings"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
)

type ValidateAccessNormalized struct {
	Token string
}

func (m *AuthMapper) ValidateAccessInputToNormalized(input dto.ValidateAccessCommand) (ValidateAccessNormalized, error) {
	vs := m.validator.ValidateStruct(input)
	token := strings.TrimSpace(strings.ToLower(input.AccessToken))
	if after, ok := strings.CutPrefix(token, "bearer "); ok {
		token = after
	}
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return ValidateAccessNormalized{}, err
	}
	return ValidateAccessNormalized{Token: token}, nil
}
