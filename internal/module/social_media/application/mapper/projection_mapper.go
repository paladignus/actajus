// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/social_media/application/readmodel"
	"github.com/paladignus/actajus/internal/module/social_media/domain"
)

type SocialMediaProjectionMapper struct{}

func NewSocialMediaProjectionMapper() SocialMediaProjectionMapper {
	return SocialMediaProjectionMapper{}
}

func (s SocialMediaProjectionMapper) ProjectSocialMediaToReadModel(socialMedia *domain.SocialMedia) *readmodel.SocialMediaReadModel {
	return &readmodel.SocialMediaReadModel{
		ID:        socialMedia.ID().Value(),
		Platform:  socialMedia.Platform().Value(),
		URL:       socialMedia.URL().Value(),
		CreatedAt: socialMedia.CreatedAt(),
		UpdatedAt: socialMedia.UpdatedAt(),
	}
}
