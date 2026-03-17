// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/address/application/dto"
	addressv1 "github.com/paladignus/actajus/proto/address/v1"
)

func ProtoToAddressCreateCommand(in *addressv1.CreateAddressRequest) dto.CreateAddressRequest {
	if in == nil {
		return dto.CreateAddressRequest{}
	}
	return dto.CreateAddressRequest{
		ZIP:          in.Zip,
		Title:        in.Title,
		Street:       in.Street,
		Number:       uint(in.Number),
		Complement:   in.Complement,
		Reference:    in.Reference,
		Neighborhood: in.Neighborhood,
		City:         in.City,
		State:        in.State,
		Country:      in.Country,
	}
}
