// Package repository
package repository

import "context"

type PermissionChecker interface {
	HasPermission(ctx context.Context, idUser int64, perm string) (bool, error)
}
