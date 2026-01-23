// Package domain
package domain

import (
	"context"
)

type SocialMediaRepository interface {
	Create(ctx context.Context, socialMedia *SocialMedia) error
}
