// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/social_media/application/dto"
	socialmediav1 "github.com/paladignus/actajus/proto/social_media/v1"
)

func ProtoToSocialMediaCreateCommands(in []*socialmediav1.CreateSocialMediaRequest) []dto.CreateSocialMediaRequest {
	if in == nil {
		return []dto.CreateSocialMediaRequest{}
	}
	result := make([]dto.CreateSocialMediaRequest, len(in))
	for i := range in {
		result[i] = dto.CreateSocialMediaRequest{
			Platform: in[i].Platform,
			URL:      in[i].Url,
		}
	}
	return result
}

func ProtoToSocialMediaUpdateCommands(in []*socialmediav1.UpdateSocialMediaRequest) []dto.UpdateSocialMediaRequest {
	if in == nil {
		return []dto.UpdateSocialMediaRequest{}
	}
	result := make([]dto.UpdateSocialMediaRequest, len(in))
	for i := range in {
		result[i] = dto.UpdateSocialMediaRequest{
			IDSocialMedia: in[i].Id,
			Platform:      in[i].Platform,
			URL:           in[i].Url,
		}
	}
	return result
}
