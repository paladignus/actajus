// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/social_media/application/dto"
	socialmediav1 "github.com/paladignus/actajus/proto/social_media/v1"
)

func ProtoToSocialMediaCreateCommands(in *socialmediav1.CreateSocialMediaRequest) []dto.CreateSocialMediaRequest {
	if in == nil || len(in.Platform) == 0 {
		return []dto.CreateSocialMediaRequest{}
	}
	result := make([]dto.CreateSocialMediaRequest, len(in.Platform))
	for i := range in.Platform {
		result[i] = dto.CreateSocialMediaRequest{
			Platform: in.Platform[i],
			URL:      in.Url[i],
		}
	}
	return result
}
