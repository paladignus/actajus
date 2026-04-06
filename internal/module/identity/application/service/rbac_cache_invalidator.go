// Package service
package service

import (
	"context"
)

type RBACCacheInvalidator interface {
	InvalidateUser(ctx context.Context, uid int64)
}
