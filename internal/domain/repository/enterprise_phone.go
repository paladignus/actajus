// Package repository
package repository

import (
	"context"
)

type IEnterprisePhone interface {
	Create(ctx context.Context, idEnterprise, idPhone uint) error
}
