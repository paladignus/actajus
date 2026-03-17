// Package domain
package domain

import "context"

type CompanyAddressRepository interface {
	Create(ctx context.Context, idCompany, idAddress int64) error
	DeleteByIDCompany(ctx context.Context, idCompany int64) error
}
