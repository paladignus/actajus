// Package adapter
package adapter

import (
	"context"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/internal/shared/infrastructure/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type SMTPEmail struct {
	config config.SMTPConfig
}

func NewSMTPEmail(config config.SMTPConfig) SMTPEmail {
	return SMTPEmail{config}
}

func (s SMTPEmail) SendEmail(ctx context.Context, to, subject, body string) error {
	from := mail.Address{Name: "Actajus", Address: s.config.From}
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", from.String(), to, subject, body)
	auth := smtp.PlainAuth("", s.config.User, s.config.Pass, s.config.Host)
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	err := smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(message))
	if err == nil {
		metrics.EmailSentCount.With(prometheus.Labels{"template": "generic"}).Inc()
	}
	return err
}
