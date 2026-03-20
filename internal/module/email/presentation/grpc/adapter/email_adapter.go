// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/email/application/command"
	"github.com/paladignus/actajus/internal/module/email/application/readmodel"
	emailv1 "github.com/paladignus/actajus/proto/email/v1"
)

func ProtoToEmailCreateCommand(in *emailv1.CreateEmailRequest) command.CreateEmailCommand {
	if in == nil {
		return command.CreateEmailCommand{}
	}
	return command.CreateEmailCommand{
		Address: in.Address,
	}
}

func ProtoToEmailUpdateCommand(in *emailv1.UpdateEmailRequest) command.UpdateEmailCommand {
	if in == nil {
		return command.UpdateEmailCommand{}
	}
	return command.UpdateEmailCommand{
		IDEmail: in.Id,
		Address: in.Address,
	}
}

func EmailReadModelToProto(in []*readmodel.EmailReadModel) []*emailv1.EmailResponse {
	emails := make([]*emailv1.EmailResponse, len(in))
	for i := range in {
		emails[i] = &emailv1.EmailResponse{
			Id:      int64(in[i].ID),
			Address: in[i].Address,
		}
	}
	return emails
}
