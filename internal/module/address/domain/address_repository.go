// Package domain
package domain

import "context"

type AddressRepository interface {
	FindByCEP(ctx context.Context, cep string) (*Address, error)
	Create(ctx context.Context, address *Address) error
}
