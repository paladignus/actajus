// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/email/application/command"
	"github.com/paladignus/actajus/internal/module/email/domain"
)

type EmailMapper struct{}

func NewEmailMapper() *EmailMapper {
	return &EmailMapper{}
}

func (e *EmailMapper) InputToDomain(input command.CreateEmailCommand) (*domain.Email, error) {
	return domain.NewEmailBuilder().
		WithAddress(input.Address).
		Build()
}

func (e *EmailMapper) UpdateInputToDomain(input command.UpdateEmailCommand) (*domain.Email, error) {
	return domain.NewEmailBuilder().
		WithID(input.IDEmail).
		WithAddress(input.Address).
		WithUpdatedAt(time.Now()).
		Build()
}
