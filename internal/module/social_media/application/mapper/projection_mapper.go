// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/social_media/application/dto"
	"github.com/paladignus/actajus/internal/module/social_media/domain"
)

type SocialMediaProjectionMapper struct{}

func NewSocialMediaProjectionMapper() SocialMediaProjectionMapper {
	return SocialMediaProjectionMapper{}
}

func (s SocialMediaProjectionMapper) ProjectSocialMediaToReadModel(socialMedia *domain.SocialMedia) *dto.SocialMediaReadModel {
	return &dto.SocialMediaReadModel{
		ID:        socialMedia.ID(),
		Platform:  socialMedia.Platform().Value(),
		URL:       socialMedia.URL().Value(),
		CreatedAt: socialMedia.CreatedAt().Format(time.RFC3339),
		UpdatedAt: socialMedia.UpdatedAt().Format(time.RFC3339),
	}
}
