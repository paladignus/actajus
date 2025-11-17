// Package adapter
package adapter

import (
	"context"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/paladignus/actajus/internal/infrastructure/config"
)

type SMTPEmail struct {
	config config.SMTPConfig
}

func NewSMTPEmail(config config.SMTPConfig) SMTPEmail {
	return SMTPEmail{config}
}

func (s SMTPEmail) SendEmail(ctx context.Context, to string, resetURL string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
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
	return smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(message))
}
