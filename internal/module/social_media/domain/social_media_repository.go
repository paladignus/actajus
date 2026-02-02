// Package domain
package domain

import (
	"context"
)

type SocialMediaRepository interface {
	Create(ctx context.Context, socialMedia *SocialMedia) error
	Update(ctx context.Context, socialMedia SocialMedia) error
	Delete(ctx context.Context, socialMedia SocialMedia) error
}
