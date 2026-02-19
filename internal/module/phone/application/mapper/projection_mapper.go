// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/phone/application/dto"
	"github.com/paladignus/actajus/internal/module/phone/domain"
)

type PhoneProjectionMapper struct{}

func NewPhoneProjectionMapper() PhoneProjectionMapper {
	return PhoneProjectionMapper{}
}

func (p PhoneProjectionMapper) ProjectPhoneToReadModel(phone *domain.Phone) *dto.PhoneReadModel {
	return &dto.PhoneReadModel{
		ID:         phone.ID(),
		Number:     phone.Number().Value(),
		Kind:       phone.Kind().Value(),
		Department: phone.Department().Value(),
		CreatedAt:  phone.CreatedAt(),
		UpdatedAt:  phone.UpdatedAt(),
	}
}
