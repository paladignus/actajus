// Package repository
package repository

import "context"

type IAddressEnterprise interface {
	Create(ctx context.Context, idEnterprise, idAddress uint) error
}
