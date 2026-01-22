// Package mapper
package mapper

import (
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
