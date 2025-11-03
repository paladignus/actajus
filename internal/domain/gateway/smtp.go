// Package gateway
package gateway

import "context"

type SMTP interface {
	SendEmail(context.Context, string, string) error
}
