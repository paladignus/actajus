// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/social_media/application/dto"
	"github.com/paladignus/actajus/internal/module/social_media/domain"
)

type SocialMediaMapper struct{}

func NewSocialMediaMapper() *SocialMediaMapper {
	return &SocialMediaMapper{}
}

func (s *SocialMediaMapper) InputToDomain(input dto.CreateSocialMediaRequest) (*domain.SocialMedia, error) {
	return domain.NewSocialMediaBuilder().
		WithName(input.Name).
		WithURL(input.URL).
		Build()
}
