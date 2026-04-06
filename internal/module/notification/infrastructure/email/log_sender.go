// Package email
package email

import (
	"context"
	"fmt"
)

type LogSender struct{}

func NewLogSender() *LogSender { return &LogSender{} }

func (s *LogSender) SendTemplate(ctx context.Context, to string, template string, data map[string]any) error {
	fmt.Printf("send email -> to=%s template=%s data=%v\n", to, template, data)
	return nil
}
