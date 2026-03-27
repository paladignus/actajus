// Package repository
package repository

import (
	"context"

	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
)

type PhoneRepository interface {
	Create(ctx context.Context, phone *phoneDomain.Phone) error
	Update(ctx context.Context, phone phoneDomain.Phone) error
	Delete(ctx context.Context, phone phoneDomain.Phone) error
	FindByIDCompany(ctx context.Context, idCompany int64) (*phoneDomain.Phone, error)
}
