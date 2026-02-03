// Package domain
package domain

import (
	"context"
)

type EmailRepository interface {
	Create(ctx context.Context, email *Email) error
	Update(ctx context.Context, email Email) error
	Delete(ctx context.Context, email Email) error
	FindByIDCompany(ctx context.Context, idCompany uint) (*Email, error)
}
