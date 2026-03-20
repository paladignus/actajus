// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/email/application/readmodel"
	"github.com/paladignus/actajus/internal/module/email/domain"
)

type EmailProjectionMapper struct{}

func NewEmailProjectionMapper() EmailProjectionMapper {
	return EmailProjectionMapper{}
}

func (e EmailProjectionMapper) ProjectEmailToReadModel(email *domain.Email) *readmodel.EmailReadModel {
	return &readmodel.EmailReadModel{
		ID:        email.ID().Value(),
		Address:   email.Address().Value(),
		CreatedAt: email.CreatedAt(),
		UpdatedAt: email.UpdatedAt(),
	}
}
