// Package repository
package repository

import (
	"context"

	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
)

type SocialMediaRepository interface {
	Create(ctx context.Context, socialMedia *socialMediaDomain.SocialMedia) error
	Update(ctx context.Context, socialMedia socialMediaDomain.SocialMedia) error
	Delete(ctx context.Context, socialMedia socialMediaDomain.SocialMedia) error
	FindByIDCompany(ctx context.Context, idCompany int64) ([]*socialMediaDomain.SocialMedia, error)
	DeleteByIDCompany(ctx context.Context, idCompany int64) error
}
