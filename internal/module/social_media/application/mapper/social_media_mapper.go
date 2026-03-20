// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/social_media/application/command"
	"github.com/paladignus/actajus/internal/module/social_media/domain"
)

type SocialMediaMapper struct{}

func NewSocialMediaMapper() *SocialMediaMapper {
	return &SocialMediaMapper{}
}

func (s *SocialMediaMapper) InputToDomain(input command.CreateSocialMediaCommand) (*domain.SocialMedia, error) {
	return domain.NewSocialMediaBuilder().
		WithPlatform(input.Platform).
		WithURL(input.URL).
		Build()
}

func (s *SocialMediaMapper) UpdateInputToDomain(input command.UpdateSocialMediaCommand) (*domain.SocialMedia, error) {
	return domain.NewSocialMediaBuilder().
		WithID(input.IDSocialMedia).
		WithPlatform(input.Platform).
		WithURL(input.URL).
		WithUpdatedAt(time.Now()).
		Build()
}
