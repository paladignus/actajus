// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/address/application/dto"
	"github.com/paladignus/actajus/internal/module/address/domain"
)

type AddressPrejectionMapper struct{}

func NewAddressProjectionMapper() AddressPrejectionMapper {
	return AddressPrejectionMapper{}
}

func (a AddressPrejectionMapper) ProjectAddressToReadModel(address *domain.Address) *dto.AddressReadModel {
	return &dto.AddressReadModel{
		ID:           address.ID(),
		ZIP:          address.ZIP().Value(),
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
