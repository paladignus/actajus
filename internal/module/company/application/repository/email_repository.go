// Package repository
package repository

import (
	"context"

	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
)

type EmailRepository interface {
	Create(ctx context.Context, email *emailDomain.Email) error
	Update(ctx context.Context, email emailDomain.Email) error
	Delete(ctx context.Context, email emailDomain.Email) error
	FindByIDCompany(ctx context.Context, idCompany int64) (*emailDomain.Email, error)
}
