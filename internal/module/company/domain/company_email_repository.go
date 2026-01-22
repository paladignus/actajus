// Package domain
package domain

import "context"

type CompanyEmailRepository interface {
	Create(ctx context.Context, idCompany, idEmail uint) error
	// DeleteByIDCompany(ctx context.Context, idCompany uint) error
	// DeleteByIDAddress(ctx context.Context, idAddress uint) error
}
