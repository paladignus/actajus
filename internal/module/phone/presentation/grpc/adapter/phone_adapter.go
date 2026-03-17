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
