// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/phone/application/dto"
	"github.com/paladignus/actajus/internal/module/phone/domain"
)

type PhoneMapper struct{}

func NewPhoneMapper() *PhoneMapper {
	return &PhoneMapper{}
}

func (p *PhoneMapper) InputToDomain(input dto.CreatePhoneRequest) (*domain.Phone, error) {
	return domain.NewPhoneBuilder().
		WithNumber(input.Number).
		WithKind(input.Kind).
		WithDepartment(input.Department).
		Build()
}

func (p *PhoneMapper) UpdateInputToDomain(input dto.UpdatePhoneRequest) (*domain.Phone, error) {
	return domain.NewPhoneBuilder().
		WithID(input.IDPhone).
		WithNumber(input.Number).
		WithKind(input.Kind).
		WithDepartment(input.Department).
		WithUpdatedAt(time.Now()).
		Build()
}
