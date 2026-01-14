// Package repository
package repository

import "context"

type IEmailEnterprise interface {
	Create(ctx context.Context, idEmail, idEnterprise uint) error
}
