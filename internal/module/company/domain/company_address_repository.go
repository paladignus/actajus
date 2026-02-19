// Package domain
package domain

import "context"

type CompanyAddressRepository interface {
	Create(ctx context.Context, idCompany, idAddress uint) error
	DeleteByIDCompany(ctx context.Context, idCompany uint) error
}
