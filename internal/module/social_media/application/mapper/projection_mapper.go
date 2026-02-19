// Package mapper
package mapper

import (
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
		CreatedAt: socialMedia.CreatedAt(),
		UpdatedAt: socialMedia.UpdatedAt(),
	}
}
