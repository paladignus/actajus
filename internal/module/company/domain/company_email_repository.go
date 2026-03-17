// Package domain
package domain

import "context"

type CompanyEmailRepository interface {
	Create(ctx context.Context, idCompany, idEmail int64) error
	DeleteByIDCompany(ctx context.Context, idCompany int64) error
}
