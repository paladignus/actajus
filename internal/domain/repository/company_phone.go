// Package repository
package repository

import (
	"context"
)

type ICompanyPhone interface {
	Create(ctx context.Context, idCompany, idPhone uint) error
}
