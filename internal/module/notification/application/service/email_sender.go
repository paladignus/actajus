// Package service
package service

import "context"

type EmailSender interface {
	SendTemplate(ctx context.Context, to string, template string, data map[string]any) error
}
