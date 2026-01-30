// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/social_media/application/dto"
	"github.com/paladignus/actajus/internal/module/social_media/domain"
)

type SocialMediaMapper struct{}

func NewSocialMediaMapper() *SocialMediaMapper {
	return &SocialMediaMapper{}
}

func (s *SocialMediaMapper) InputToDomain(input dto.CreateSocialMediaRequest) (*domain.SocialMedia, error) {
	return domain.NewSocialMediaBuilder().
		WithPlatform(input.Platform).
		WithURL(input.URL).
		Build()
}

func (s *SocialMediaMapper) UpdateInputToDomain(input dto.UpdateSocialMediaRequest) (*domain.SocialMedia, error) {
	return domain.NewSocialMediaBuilder().
		WithID(input.IDSocialMedia).
		WithPlatform(input.Platform).
		WithURL(input.URL).
		WithUpdatedAt(time.Now()).
		Build()
}
