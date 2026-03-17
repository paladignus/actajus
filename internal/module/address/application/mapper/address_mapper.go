// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/address/application/dto"
	"github.com/paladignus/actajus/internal/module/address/domain"
)

type AddressMapper struct{}

func NewAddressMapper() *AddressMapper {
	return &AddressMapper{}
}

func (m *AddressMapper) InputToDomain(input dto.CreateAddressRequest) (*domain.Address, error) {
	return domain.NewAddressBuilder().
		WithZIP(input.ZIP).
		WithTitle(input.Title).
		WithStreet(input.Street).
		WithNumber(input.Number).
		WithComplement(input.Complement).
		WithReference(input.Reference).
		WithNeighborhood(input.Neighborhood).
		WithCity(input.City).
		WithState(input.State).
		WithCountry(input.Country).
		Build()
}

func (m *AddressMapper) UpdateInputToDomain(input dto.UpdateAddressRequest) (*domain.Address, error) {
	return domain.NewAddressBuilder().
		WithID(input.IDAddress).
		WithZIP(input.ZIP).
		WithTitle(input.Title).
		WithStreet(input.Street).
		WithNumber(input.Number).
		WithComplement(input.Complement).
		WithReference(input.Reference).
		WithNeighborhood(input.Neighborhood).
		WithCity(input.City).
		WithState(input.State).
		WithCountry(input.Country).
		WithUpdatedAt(time.Now()).
		Build()
}
