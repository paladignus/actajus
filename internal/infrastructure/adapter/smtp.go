// Package adapter
package adapter

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/paladignus/actajus/internal/infrastructure/config"
)

type SMTPEmail struct {
	config config.SMTPConfig
}

func NewSMTPEmail(config config.SMTPConfig) SMTPEmail {
	return SMTPEmail{config}
}

func (s *SMTPEmail) SendEmail(ctx context.Context, to string, resetURL string) error {
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
        		<h2>Password Reset Request</h2>
        		<p>You have requested to reset your password. Click the button below to proceed:</p>
        		<a href="%s" class="button">Reset Password</a>
        		<p>Or copy and paste this link into your browser:</p>
        		<p><a href="%s">%s</a></p>
        		<p><strong>This link will expire in 30 minutes.</strong></p>
        		<p>If you did not request this password reset, please ignore this email.</p>
        		<div class="footer">
            	<p>This is an automated message, please do not reply.</p>
        		</div>
    			</div>
				</body>
		</html>
	`, resetURL, resetURL, resetURL)
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.config.From, to, subject, body)
	auth := smtp.PlainAuth("", s.config.User, s.config.Pass, s.config.Host)
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	return smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(message))
}
