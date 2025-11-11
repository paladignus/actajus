// Package event
package event

import (
	"context"
	"fmt"
	"log"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/gateway"
)

// SendEmailHandler processa eventos de domínio que requerem envio de email.
// Especificamente, escuta o evento "user.password_reset_requested" e
// envia o email de recuperação de senha.
type SendEmailHandler struct {
	// smtpGateway é o adapter que envia emails (sua implementação existente)
	smtpGateway gateway.SMTP
}

// Garante em tempo de compilação que SendEmailHandler implementa event.Handler
var _ event.Handler = (*SendEmailHandler)(nil)

// NewSendEmailHandler cria um novo handler para envio de emails.
//
// Parâmetros:
//   - smtpGateway: implementação do gateway SMTP (seu adapter existente)
//
// Retorna:
//   - *SendEmailHandler: handler configurado e pronto para uso
func NewSendEmailHandler(smtpGateway gateway.SMTP) *SendEmailHandler {
	return &SendEmailHandler{
		smtpGateway: smtpGateway,
	}
}

// Handle processa o evento de solicitação de reset de senha.
// É chamado automaticamente pelo Subscriber quando uma mensagem chega.
//
// Parâmetros:
//   - ctx: contexto com timeout/cancelamento
//   - evt: o evento recebido
//
// Retorna:
//   - error: se falhar ao enviar email
//
// Importante: Este método DEVE ser idempotente!
// Se chamado múltiplas vezes com o mesmo evento, deve produzir
// o mesmo resultado sem efeitos colaterais indesejados.
func (h *SendEmailHandler) Handle(ctx context.Context, evt event.Event) error {
	// Log para observabilidade
	log.Printf("[SendEmailHandler] Processing event: %s (aggregate: %s)",
		evt.EventName(), evt.GetAggregateID())

	// Type assertion para o tipo específico (agora type-safe graças ao Registry!)
	resetEvt, ok := evt.(*event.PasswordResetRequestedEvent)
	if !ok {
		return fmt.Errorf("expected *PasswordResetRequestedEvent, got %T", evt)
	}

	// Extrai dados do evento (agora com acesso direto aos campos!)
	userEmail := resetEvt.UserEmail
	resetURL := resetEvt.ResetURL
	// userName := resetEvt.UserName

	// Prepara o assunto do email
	// subject := "Requisição de Recuperação de Senha"

	// Monta o corpo do email
	// body := h.buildEmailBody(userName, resetURL)

	// Envia o email usando seu adapter SMTP existente
	log.Printf("[SendEmailHandler] Sending email to: %s", userEmail)

	err := h.smtpGateway.SendEmail(ctx, userEmail, resetURL)
	if err != nil {
		// Log do erro para troubleshooting
		log.Printf("[SendEmailHandler] Failed to send email to %s: %v", userEmail, err)

		// Retorna erro - o NATS fará retry automático
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("[SendEmailHandler] Email sent successfully to: %s", userEmail)
	return nil
}

// CanHandle verifica se este handler pode processar o evento dado.
// Retorna true apenas para eventos de reset de senha.
//
// Parâmetros:
//   - evt: o evento a verificar
//
// Retorna:
//   - true se for "user.password_reset_requested"
//   - false para qualquer outro evento
func (h *SendEmailHandler) CanHandle(evt event.Event) bool {
	return evt.EventName() == "user.password_reset_requested"
}

// parseEvent converte o evento genérico para um map com os dados específicos.
// Você pode melhorar isso criando types específicos para cada evento.
// func (h *SendEmailHandler) parseEvent(evt event.Event) (map[string]any, error) {
// 	// Serializa e desserializa para obter todos os campos
// 	// (workaround simples - em produção, use type assertion ou reflection)
// 	data, err := json.Marshal(evt)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var result map[string]any
// 	if err := json.Unmarshal(data, &result); err != nil {
// 		return nil, err
// 	}
//
// 	// Valida campos obrigatórios
// 	if _, ok := result["UserEmail"]; !ok {
// 		return nil, fmt.Errorf("missing required field: UserEmail")
// 	}
// 	if _, ok := result["ResetURL"]; !ok {
// 		return nil, fmt.Errorf("missing required field: ResetURL")
// 	}
//
// 	return result, nil
// }

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

// Exemplo de como criar handlers adicionais no futuro:
//
// type WelcomeEmailHandler struct {
//     smtpGateway gateway.SMTP
// }
//
// func (h *WelcomeEmailHandler) Handle(ctx context.Context, evt event.Event) error {
//     // Envia email de boas-vindas
// }
//
// func (h *WelcomeEmailHandler) CanHandle(evt event.Event) bool {
//     return evt.EventName() == "user.created"
// }

// type SendEmailHandler struct {
// 	smtpGateway gateway.SMTP
// }
//
// func NewSendEmailHandler(smtpGateway gateway.SMTP) *SendEmailHandler {
// 	return &SendEmailHandler{
// 		smtpGateway: smtpGateway,
// 	}
// }
//
// func (h *SendEmailHandler) Handle(ctx context.Context, evt event.Event) error {
// 	log.Printf("[SendEmailHandler] Processing event: %s (aggregate: %s)", evt.EventName(), evt.GetAggregateID())
// 	var passwordResetEvt event.PasswordResetRequestedEvent
// 	if ptrEvt, ok := evt.(*event.PasswordResetRequestedEvent); ok {
// 		passwordResetEvt = *ptrEvt // Desreferencia o ponteiro para obter o valor
// 	} else if valEvt, ok := evt.(event.PasswordResetRequestedEvent); ok {
// 		passwordResetEvt = valEvt
// 	} else {
// 		return fmt.Errorf("handler received unexpected event type: %T, expected *event.PasswordResetRequestedEvent or event.PasswordResetRequestedEvent", evt)
// 	}
//
// 	userEmail := passwordResetEvt.UserEmail
// 	resetURL := passwordResetEvt.ResetURL
// 	if userEmail == "" || resetURL == "" {
// 		return fmt.Errorf("missing required fields in PasswordResetRequestedEvent: UserEmail='%s', ResetURL='%s'", userEmail, resetURL)
// 	}
// 	log.Printf("[SendEmailHandler] Sending email to: %s", userEmail)
// 	err := h.smtpGateway.SendEmail(ctx, userEmail, resetURL)
// 	if err != nil {
// 		log.Printf("[SendEmailHandler] Failed to send email to %s: %v", userEmail, err)
// 		return fmt.Errorf("failed to send email: %w", err)
// 	}
// 	log.Printf("[NATS] Email sent successfully to: %s", userEmail)
// 	return nil
// }
//
// func (h *SendEmailHandler) CanHandle(evt event.Event) bool {
// 	return evt.EventName() == "user.password_reset_requested"
// }
