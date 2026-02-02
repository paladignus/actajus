// Package domain
package domain

import "context"

type AddressRepository interface {
	Create(ctx context.Context, address *Address) error
	Update(ctx context.Context, address Address) error
	// FindByCEP(ctx context.Context, cep string) (Address, error)
	FindByIDCompany(ctx context.Context, idCompany uint) (*Address, error)
}
