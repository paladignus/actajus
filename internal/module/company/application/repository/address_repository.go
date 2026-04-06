// Package repository
package repository

import (
	"context"

	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
)

type AddressRepository interface {
	Create(ctx context.Context, address *addrDomain.Address) error
	Update(ctx context.Context, address addrDomain.Address) error
	FindByIDCompany(ctx context.Context, idCompany int64) (*addrDomain.Address, error)
}
