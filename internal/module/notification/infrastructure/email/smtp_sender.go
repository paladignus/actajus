// Package email
package email

import (
	"context"
	"fmt"
	"net/mail"
	"net/smtp"
)

type SMTPSender struct {
	host     string
	port     int
	user     string
	pass     string
	from     string
	fromName string
}

func NewSMTPSender(host string, port int, user, pass, from string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		user:     user,
		pass:     pass,
		from:     from,
		fromName: "Actajus",
	}
}

func (s *SMTPSender) SendHTML(ctx context.Context, to string, subject string, htmlBody string) error {
	from := mail.Address{Name: s.fromName, Address: s.from}
	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		from.String(),
		to,
		subject,
		htmlBody,
	)
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
