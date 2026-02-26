// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/domain"
)

type Authorization interface {
	ListPermissionsByUser(ctx context.Context, idUser domain.IDUser) ([]string, error)
}
