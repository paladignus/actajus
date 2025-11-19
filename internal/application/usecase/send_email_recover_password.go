// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type SendEmailRecoverPassword struct {
	logger     repository.Logger
	smtp       gateway.SMTP
	subscriber gateway.Subscriber
}

func NewSendEmailRecoverPassword(
	repository repository.Logger,
	smtp gateway.SMTP,
	subscriber gateway.Subscriber,
) SendEmailRecoverPassword {
	return SendEmailRecoverPassword{
		repository,
		smtp,
		subscriber,
	}
}

func (s SendEmailRecoverPassword) Execute(ctx context.Context) error {
	s.logger.Info(ctx, "send email recover password")
	eventName := "user.password_reset_requested"
	handler := func(ctx context.Context, evt event.IEvent) error {
		resetEvt, ok := evt.(*event.PasswordResetRequestedEvent)
		if !ok {
			return fmt.Errorf("expected *PasswordResetRequestedEvent, got %T", evt)
		}
		userEmail := resetEvt.UserEmail
		resetURL := resetEvt.ResetURL
		s.logger.Info(ctx, "sending email", "to", userEmail)
		err := s.smtp.SendEmail(ctx, userEmail, resetURL)
		if err != nil {
			s.logger.Error(ctx, "failed to send email", "to", userEmail, "error", err)
			return fmt.Errorf("failed to send email: %w", err)
		}
		s.logger.Info(ctx, "email sent successfully", "to", userEmail)
		return nil
	}
	err := s.subscriber.Subscribe(ctx, eventName, handler)
	if err != nil {
		s.logger.Error(ctx, "failed to subscribe to event", "event", eventName, "error", err)
		return fmt.Errorf("failed to subscribe handler for %s: %w", eventName, err)
	}
	s.logger.Info(ctx, "handler registered for event", "event", eventName)
	return nil
}
