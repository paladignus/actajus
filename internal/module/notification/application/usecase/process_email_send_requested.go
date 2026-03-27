// Package usecase
package usecase

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/paladignus/actajus/internal/module/notification/application/message"
	"github.com/paladignus/actajus/internal/module/notification/application/repository"
	"github.com/paladignus/actajus/internal/module/notification/application/service"
)

type ProcessEmailSendRequested struct {
	sender     service.EmailSender
	renderer   service.TemplateRenderer
	sent       repository.SentEmailRepository
	publicBase string
}

func NewProcessEmailSendRequested(
	sender service.EmailSender,
	renderer service.TemplateRenderer,
	sent repository.SentEmailRepository,
	publicBase string,
) ProcessEmailSendRequested {
	return ProcessEmailSendRequested{
		sender,
		renderer,
		sent,
		publicBase,
	}
}

func (uc ProcessEmailSendRequested) Execute(ctx context.Context, msg message.EmailSendRequested) error {
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
	if msg.SentAfter != nil && time.Now().Before(*msg.SentAfter) {
		return fmt.Errorf("not time yet")
	}
	if msg.Template == "password-reset" {
		resetURL, err := uc.buildResetURL(msg.Data)
		if err != nil {
			return fmt.Errorf("build reset url: %w", err)
		}
		msg.Data["reset_url"] = resetURL
	}
	subject, html, err := uc.renderer.RenderHTML(ctx, msg.Template, msg.Data)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	if err := uc.sender.SendHTML(ctx, msg.To, subject, html); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	if err := uc.sent.MarkSent(ctx, msg.IDMessage, msg.To, msg.Template); err != nil {
		return fmt.Errorf("mark sent: %w", err)
	}
	return nil
}

func (uc ProcessEmailSendRequested) buildResetURL(data map[string]any) (string, error) {
	base, err := url.Parse(uc.publicBase)
	if err != nil {
		return "", fmt.Errorf("invalid public base url: %w", err)
	}
	base.Path = "/reset-password"
	q := base.Query()
	idReset, ok := data["id_reset"]
	if !ok {
		return "", fmt.Errorf("missing data.id_reset")
	}
	token, ok := data["token"]
	if !ok {
		return "", fmt.Errorf("missing data.token")
	}
	q.Set("id_reset", fmt.Sprintf("%v", idReset))
	q.Set("token", fmt.Sprintf("%v", token))
	base.RawQuery = q.Encode()
	return base.String(), nil
}
