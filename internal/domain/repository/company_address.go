// Package repository
package repository

import "context"

type ICompanyAddress interface {
	Create(ctx context.Context, idCompany, idAddress uint) error
}
