// Package service
package service

import "context"

type EmailSender interface {
	SendHTML(ctx context.Context, to string, subject string, body string) error
}
