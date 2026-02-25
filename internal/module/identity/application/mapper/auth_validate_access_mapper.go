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
	token := strings.TrimSpace(input.AccessToken)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = strings.TrimSpace(token[7:])
	}
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return ValidateAccessNormalized{}, err
	}
	return ValidateAccessNormalized{Token: token}, nil
}
