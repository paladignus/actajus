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
	// if err != nil {
	// 	return nil, err
	// }
	// return address, nil
}

func (m *AddressMapper) DomainToOutput(address *domain.Address) dto.AddressResponse {
	return dto.AddressResponse{
		ID:           address.ID(),
		ZIP:          address.ZIP().Formatted(),
		Title:        address.Title().Value(),
		Street:       address.Street().Value(),
		Number:       address.Number(),
		Complement:   address.Complement().Value(),
		Reference:    address.Reference().Value(),
		Neighborhood: address.Neighborhood().Value(),
		City:         address.City().Value(),
		State:        address.State().Value(),
		Country:      address.Country().Value(),
		CreatedAt:    address.CreatedAt().Format(time.RFC3339),
		UpdatedAt:    address.UpdatedAt().Format(time.RFC3339),
	}
}
