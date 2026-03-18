// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/phone/application/dto"
	phonev1 "github.com/paladignus/actajus/proto/phone/v1"
)

func ProtoToPhoneCreateCommand(in *phonev1.CreatePhoneRequest) dto.CreatePhoneRequest {
	if in == nil {
		return dto.CreatePhoneRequest{}
	}
	return dto.CreatePhoneRequest{
		Number:     in.Number,
		Kind:       in.Kind,
		Department: in.Department,
	}
}

func ProtoToPhoneUpdateCommand(in *phonev1.UpdatePhoneRequest) dto.UpdatePhoneRequest {
	if in == nil {
		return dto.UpdatePhoneRequest{}
	}
	return dto.UpdatePhoneRequest{
		IDPhone:    in.Id,
		Number:     in.Number,
		Kind:       in.Kind,
		Department: in.Department,
	}
}

// func AddressReadModelToProto(in *dto.AddressReadModel) *addressv1.AddressResponse {

func PhoneReadModelToProto(in []*dto.PhoneReadModel) []*phonev1.PhoneResponse {
	phones := make([]*phonev1.PhoneResponse, len(in))
	for i := range in {
		phones[i] = &phonev1.PhoneResponse{
			Id:         int64(in[i].ID),
			Number:     in[i].Number,
			Kind:       in[i].Kind,
			Department: in[i].Department,
		}
	}
	return phones
}
