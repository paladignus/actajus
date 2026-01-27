// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/email/application/dto"
	"github.com/paladignus/actajus/internal/module/email/domain"
)

type EmailProjectionMapper struct{}

func NewEmailProjectionMapper() EmailProjectionMapper {
	return EmailProjectionMapper{}
}

func (e EmailProjectionMapper) ProjectEmailToReadModel(email *domain.Email) *dto.EmailReadModel {
	return &dto.EmailReadModel{
		ID:        email.ID(),
		Address:   email.Address().Value(),
		CreatedAt: email.CreatedAt().Format(time.RFC3339),
		UpdatedAt: email.UpdatedAt().Format(time.RFC3339),
	}
}
