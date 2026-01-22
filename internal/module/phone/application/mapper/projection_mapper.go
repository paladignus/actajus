// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/phone/application/dto"
	"github.com/paladignus/actajus/internal/module/phone/domain"
)

type PhoneProjectionMapper struct{}

func NewPhoneProjectionMapper() PhoneProjectionMapper {
	return PhoneProjectionMapper{}
}

func (p PhoneProjectionMapper) ProjectPhoneToReadModel(phone *domain.Phone) *dto.PhoneReadModel {
	return &dto.PhoneReadModel{
		Number:     phone.Number().Value(),
		Kind:       phone.Kind().Value(),
		Department: phone.Department().Value(),
		CreatedAt:  phone.CreatedAt().Format(time.RFC3339),
		UpdatedAt:  phone.UpdatedAt().Format(time.RFC3339),
	}
}
