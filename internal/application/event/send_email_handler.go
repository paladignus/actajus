// Package event
package event

import (
	"context"
	"fmt"
	"log"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/gateway"
)

type SendEmailHandler struct {
	smtpGateway gateway.SMTP
}

func NewSendEmailHandler(smtpGateway gateway.SMTP) *SendEmailHandler {
	return &SendEmailHandler{
		smtpGateway: smtpGateway,
	}
}

func (h *SendEmailHandler) Handle(ctx context.Context, evt event.Event) error {
	log.Printf("[SendEmailHandler] Processing event: %s (aggregate: %s)", evt.EventName(), evt.GetAggregateID())
	var passwordResetEvt event.PasswordResetRequestedEvent
	if ptrEvt, ok := evt.(*event.PasswordResetRequestedEvent); ok {
		passwordResetEvt = *ptrEvt // Desreferencia o ponteiro para obter o valor
	} else if valEvt, ok := evt.(event.PasswordResetRequestedEvent); ok {
		passwordResetEvt = valEvt
	} else {
		return fmt.Errorf("handler received unexpected event type: %T, expected *event.PasswordResetRequestedEvent or event.PasswordResetRequestedEvent", evt)
	}

	userEmail := passwordResetEvt.UserEmail
	resetURL := passwordResetEvt.ResetURL
	if userEmail == "" || resetURL == "" {
		return fmt.Errorf("missing required fields in PasswordResetRequestedEvent: UserEmail='%s', ResetURL='%s'", userEmail, resetURL)
	}
	log.Printf("[SendEmailHandler] Sending email to: %s", userEmail)
	err := h.smtpGateway.SendEmail(ctx, userEmail, resetURL)
	if err != nil {
		log.Printf("[SendEmailHandler] Failed to send email to %s: %v", userEmail, err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	log.Printf("[NATS] Email sent successfully to: %s", userEmail)
	return nil
}

func (h *SendEmailHandler) CanHandle(evt event.Event) bool {
	return evt.EventName() == "user.password_reset_requested"
}
