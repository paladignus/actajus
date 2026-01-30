// Package mapper
package mapper

import (
	"time"

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

func (e *EmailMapper) UpdateInputToDomain(input dto.UpdateEmailRequest) (*domain.Email, error) {
	return domain.NewEmailBuilder().
		WithID(input.ID).
		WithAddress(input.Address).
		WithUpdatedAt(time.Now()).
		Build()
}
