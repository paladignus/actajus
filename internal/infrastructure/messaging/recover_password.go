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
	resetEvt, ok := evt.(*event.PasswordResetRequestedEvent)
	if !ok {
		r.logger.Error(ctx, "invalid event type", "expected", "*PasswordResetRequestedEvent", "got", fmt.Sprintf("%T", evt))
		return fmt.Errorf("expected *PasswordResetRequestedEvent, got %T", evt)
	}
	userEmail := resetEvt.UserEmail
	err := r.smtp.SendEmail(ctx, userEmail, "Recuperação de senha", body(resetEvt.ResetURL))
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

func body(resetURL string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
			<html>
				<head>
    			<meta charset="UTF-8">
    				<style>
        			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
        			.button { 
            		display: inline-block; 
            		padding: 12px 24px; 
            		background-color: #007bff; 
            		color: #ffffff !important; 
            		text-decoration: none; 
            		border-radius: 4px;
            		margin: 20px 0;
        			}
        			.footer { margin-top: 30px; font-size: 12px; color: #666; }
    				</style>
				</head>
				<body>
    			<div class="container">
        		<h2>Requisição de Recuperação de Senha</h2>
        		<p>Você solicitou a redefinição da sua senha. Clique no botão abaixo para continuar:</p>
        		<a href="%s" class="button">Recuperar Senha</a>
        		<p>Ou copie e cole este link no seu navegador:</p>
        		<p><a href="%s">%s</a></p>
        		<p><strong>Este link expirará em 30 minutos.</strong></p>
        		<p>Se você não solicitou a redefinição de senha, ignore este e-mail.</p>
        		<div class="footer">
            	<p>Esta é uma mensagem automática, por favor, não responda.</p>
        		</div>
    			</div>
				</body>
		</html>
	`, resetURL, resetURL, resetURL)
}
