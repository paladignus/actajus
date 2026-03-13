// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/notification/application/dto"
	"github.com/paladignus/actajus/internal/module/notification/application/repository"
	"github.com/paladignus/actajus/internal/module/notification/application/service"
)

type ProcessEmailSendRequested struct {
	sender service.EmailSender
	sent   repository.SentEmailRepository
}

func NewProcessEmailSendRequested(
	sender service.EmailSender,
	sent repository.SentEmailRepository,
) ProcessEmailSendRequested {
	return ProcessEmailSendRequested{
		sender,
		sent,
	}
}

func (uc ProcessEmailSendRequested) Execute(ctx context.Context, msg dto.EmailSendRequested) error {
	if msg.IDMessage == "" {
		return fmt.Errorf("message_id is required")
	}
	if msg.To == "" {
		return fmt.Errorf("to is required")
	}
	if msg.Template == "" {
		return fmt.Errorf("template is required")
	}
	exists, err := uc.sent.ExistsByIDMessage(ctx, msg.IDMessage)
	if err != nil {
		return fmt.Errorf("check idempotency: %w", err)
	}
	if exists {
		return nil
	}
	if err := uc.sender.SendTemplate(ctx, msg.To, msg.Template, msg.Data); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	if err := uc.sent.MarkSent(ctx, msg.IDMessage, msg.To, msg.Template); err != nil {
		return fmt.Errorf("mark sent: %w", err)
	}
	return nil
}
