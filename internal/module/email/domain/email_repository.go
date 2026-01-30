// Package domain
package domain

import (
	"context"
)

type EmailRepository interface {
	Create(ctx context.Context, email *Email) error
	Update(ctx context.Context, email *Email) error
}
