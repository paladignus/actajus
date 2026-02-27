// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/domain"
)

type RBACCacheInvalidator interface {
	InvalidateUser(ctx context.Context, userID domain.IDUser)
}
