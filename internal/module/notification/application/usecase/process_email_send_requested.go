// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/notification/application/dto"
	"github.com/paladignus/actajus/internal/module/notification/application/repository"
	"github.com/paladignus/actajus/internal/module/notification/application/service"
)

type ProcessEmailSendRequested struct {
	emailSender service.EmailSender
	sentEmail   repository.SentEmailRepository
}

func NewProcessEmailSendRequested(
	emailSender service.EmailSender,
	sentEmail repository.SentEmailRepository,
) ProcessEmailSendRequested {
	return ProcessEmailSendRequested{
		emailSender: emailSender,
		sentEmail:   sentEmail,
	}
}

func (uc ProcessEmailSendRequested) Execute(ctx context.Context, msg dto.EmailSendRequested) error {
	exists, err := uc.sentEmail.ExistsByIDMessage(ctx, msg.IDMessage)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if err := uc.emailSender.SendTemplate(ctx, msg.To, msg.Template, msg.Data); err != nil {
		return err
	}
	return uc.sentEmail.MarkSent(ctx, msg.IDMessage, msg.To, msg.Template)
}
