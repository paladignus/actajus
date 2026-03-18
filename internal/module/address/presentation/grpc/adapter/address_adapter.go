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

func ProtoToAddressUpdateCommand(in *addressv1.UpdateAddressRequest) dto.UpdateAddressRequest {
	if in == nil {
		return dto.UpdateAddressRequest{}
	}
	return dto.UpdateAddressRequest{
		IDAddress:    in.Id,
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

func AddressReadModelToProto(in *dto.AddressReadModel) *addressv1.AddressResponse {
	return &addressv1.AddressResponse{
		Id:           int64(in.ID),
		Zip:          in.ZIP,
		Title:        in.Title,
		Street:       in.Street,
		Complement:   *in.Complement,
		Reference:    *in.Reference,
		Number:       uint32(in.Number),
		Neighborhood: in.Neighborhood,
		City:         in.City,
		State:        in.State,
		Country:      in.Country,
	}
}
