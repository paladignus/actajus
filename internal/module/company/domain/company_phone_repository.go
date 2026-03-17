// Package domain
package domain

import "context"

type CompanyPhoneRepository interface {
	Create(ctx context.Context, idCompany, idPhone int64) error
	DeleteByIDCompany(ctx context.Context, idCompany int64) error
}
