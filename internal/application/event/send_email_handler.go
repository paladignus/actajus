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
	log.Printf("[SendEmailHandler] Processing event: %s (aggregate: %s)",
		evt.EventName(), evt.GetAggregateID())
	resetEvt, ok := evt.(*event.PasswordResetRequestedEvent)
	if !ok {
		return fmt.Errorf("expected *PasswordResetRequestedEvent, got %T", evt)
	}
	userEmail := resetEvt.UserEmail
	resetURL := resetEvt.ResetURL
	log.Printf("[SendEmailHandler] Sending email to: %s", userEmail)
	err := h.smtpGateway.SendEmail(ctx, userEmail, resetURL)
	if err != nil {
		log.Printf("[SendEmailHandler] Failed to send email to %s: %v", userEmail, err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	log.Printf("[SendEmailHandler] Email sent successfully to: %s", userEmail)
	return nil
}

func (h *SendEmailHandler) CanHandle(evt event.Event) bool {
	return evt.EventName() == "user.password_reset_requested"
}

// buildEmailBody monta o corpo do email de recuperação de senha.
// Você pode melhorar isso usando templates HTML mais elaborados.
// func (h *SendEmailHandler) buildEmailBody(userName, resetURL string) string {
// 	greeting := "Olá"
// 	if userName != "" {
// 		greeting = fmt.Sprintf("Olá, %s", userName)
// 	}
//
// 	return fmt.Sprintf(`
// %s!
//
// Você solicitou a redefinição da sua senha.
//
// Clique no botão abaixo para redefinir sua senha:
// %s
//
// Este link expirará em 30 minutos.
//
// Se você não solicitou esta redefinição, ignore este email.
//
// Atenciosamente,
// Equipe Actajus
// 	`, greeting, resetURL)
// }
