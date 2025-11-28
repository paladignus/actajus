// Package gateway
package gateway

import "context"

type SMTP interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}
