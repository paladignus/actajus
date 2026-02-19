// Package mapper
package mapper

import (
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
		CreatedAt: email.CreatedAt(),
		UpdatedAt: email.UpdatedAt(),
	}
}
