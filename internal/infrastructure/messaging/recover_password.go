// Package messaging
package messaging

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type RecoverPassword struct {
	logger repository.Logger
	smtp   gateway.SMTP
}

func NewRecoverPassword(
	logger repository.Logger,
	smtp gateway.SMTP,
) RecoverPassword {
	return RecoverPassword{
		logger, smtp,
	}
}

func (r RecoverPassword) Handle(
	ctx context.Context,
	evt event.IEvent,
) error {
	r.logger.Info(ctx, "processing", "event", evt.EventName(), "aggregate", evt.GetAggregateID())
	resetEvt, ok := evt.(*event.PasswordResetRequestedEvent)
	if !ok {
		return fmt.Errorf("expected *PasswordResetRequestedEvent, got %T", evt)
	}
	userEmail := resetEvt.UserEmail
	resetURL := resetEvt.ResetURL
	r.logger.Info(ctx, "sending email", "to", userEmail)
	err := r.smtp.SendEmail(ctx, userEmail, resetURL)
	if err != nil {
		r.logger.Error(ctx, "failed to send email", "to", userEmail, "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	r.logger.Info(ctx, "email sent successfully", "to", userEmail)
	return nil
}

func (r RecoverPassword) CanHandle(evt event.IEvent) bool {
	return evt.EventName() == "user.password_reset_requested"
}
