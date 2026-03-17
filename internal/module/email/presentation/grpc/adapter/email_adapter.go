// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/email/application/dto"
	emailv1 "github.com/paladignus/actajus/proto/email/v1"
)

func ProtoToEmailCreateCommand(in *emailv1.CreateEmailRequest) dto.CreateEmailRequest {
	if in == nil {
		return dto.CreateEmailRequest{}
	}
	return dto.CreateEmailRequest{
		Address: in.Address,
	}
}
