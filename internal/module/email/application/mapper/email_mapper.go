// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/email/application/dto"
	"github.com/paladignus/actajus/internal/module/email/domain"
)

type EmailMapper struct{}

func NewEmailMapper() *EmailMapper {
	return &EmailMapper{}
}

func (e *EmailMapper) InputToDomain(input dto.CreateEmailRequest) (*domain.Email, error) {
	return domain.NewEmailBuilder().
		WithAddress(input.Address).
		Build()
}
