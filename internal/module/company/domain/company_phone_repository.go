// Package domain
package domain

import "context"

type CompanyPhoneRepository interface {
	Create(ctx context.Context, idCompany, idPhone uint) error
	DeleteByIDCompany(ctx context.Context, idCompany uint) error
}
