// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/address/application/command"
	"github.com/paladignus/actajus/internal/module/address/application/readmodel"
	addressv1 "github.com/paladignus/actajus/proto/address/v1"
)

func ProtoToAddressCreateCommand(in *addressv1.CreateAddressRequest) command.CreateAddressCommand {
	if in == nil {
		return command.CreateAddressCommand{}
	}
	return command.CreateAddressCommand{
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

func ProtoToAddressUpdateCommand(in *addressv1.UpdateAddressRequest) command.UpdateAddressCommand {
	if in == nil {
		return command.UpdateAddressCommand{}
	}
	return command.UpdateAddressCommand{
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

func AddressReadModelToProto(in *readmodel.AddressReadModel) *addressv1.AddressResponse {
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
