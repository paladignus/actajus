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

func ProtoToEmailUpdateCommand(in *emailv1.UpdateEmailRequest) dto.UpdateEmailRequest {
	if in == nil {
		return dto.UpdateEmailRequest{}
	}
	return dto.UpdateEmailRequest{
		IDEmail: in.Id,
		Address: in.Address,
	}
}

func EmailReadModelToProto(in []*dto.EmailReadModel) []*emailv1.EmailResponse {
	// if in == nil {
	// 	return &emailv1.EmailResponse{}
	// }
	emails := make([]*emailv1.EmailResponse, len(in))
	for i := range in {
		emails[i] = &emailv1.EmailResponse{
			Id:      int64(in[i].ID),
			Address: in[i].Address,
		}
	}
	return emails
}
