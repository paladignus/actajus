// Package repository
package repository

import "context"

type SentEmailRepository interface {
	ExistsByIDMessage(ctx context.Context, mid string) (bool, error)
	MarkSent(ctx context.Context, mid string, recipient string, template string) error
}
