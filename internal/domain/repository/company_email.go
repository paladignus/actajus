// Package repository
package repository

import "context"

type ICompanyEmail interface {
	Create(ctx context.Context, idCompany, idEmail uint) error
}
