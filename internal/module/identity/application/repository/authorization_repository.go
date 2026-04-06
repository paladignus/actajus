// Package repository
package repository

import (
	"context"
)

type Authorization interface {
	ListPermissionsByUser(ctx context.Context, uid int64) ([]string, error)
}
