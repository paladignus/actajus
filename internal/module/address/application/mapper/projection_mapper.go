// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/address/application/readmodel"
	"github.com/paladignus/actajus/internal/module/address/domain"
)

type AddressPrejectionMapper struct{}

func NewAddressProjectionMapper() AddressPrejectionMapper {
	return AddressPrejectionMapper{}
}

func (a AddressPrejectionMapper) ProjectAddressToReadModel(address *domain.Address) *readmodel.AddressReadModel {
	complement := address.Complement().Value()
	reference := address.Reference().Value()
	return &readmodel.AddressReadModel{
		ID:           address.ID().Value(),
		ZIP:          address.ZIP().Value(),
		Title:        address.Title().Value(),
		Street:       address.Street().Value(),
		Number:       address.Number(),
		Complement:   &complement,
		Reference:    &reference,
		Neighborhood: address.Neighborhood().Value(),
		City:         address.City().Value(),
		State:        address.State().Value(),
		Country:      address.Country().Value(),
		CreatedAt:    address.CreatedAt(),
		UpdatedAt:    address.UpdatedAt(),
	}
}
