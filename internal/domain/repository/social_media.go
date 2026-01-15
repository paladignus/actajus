// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type ISocialMedia interface {
	Create(ctx context.Context, input entity.SocialMedia) error
}
