// Package domain
package domain

import (
	"context"
)

type PhoneRepository interface {
	Create(ctx context.Context, phone *Phone) error
	Update(ctx context.Context, phone Phone) error
	Delete(ctx context.Context, phone Phone) error
	FindByIDCompany(ctx context.Context, idCompany uint) (*Phone, error)
}
