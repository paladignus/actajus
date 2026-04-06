// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/social_media/application/command"
	"github.com/paladignus/actajus/internal/module/social_media/application/readmodel"
	socialmediav1 "github.com/paladignus/actajus/proto/social_media/v1"
)

func ProtoToSocialMediaCreateCommands(in []*socialmediav1.CreateSocialMediaRequest) []command.CreateSocialMediaCommand {
	if in == nil {
		return []command.CreateSocialMediaCommand{}
	}
	result := make([]command.CreateSocialMediaCommand, len(in))
	for i := range in {
		result[i] = command.CreateSocialMediaCommand{
			Platform: in[i].Platform,
			URL:      in[i].Url,
		}
	}
	return result
}

func ProtoToSocialMediaUpdateCommands(in []*socialmediav1.UpdateSocialMediaRequest) []command.UpdateSocialMediaCommand {
	if in == nil {
		return []command.UpdateSocialMediaCommand{}
	}
	result := make([]command.UpdateSocialMediaCommand, len(in))
	for i := range in {
		result[i] = command.UpdateSocialMediaCommand{
			IDSocialMedia: in[i].Id,
			Platform:      in[i].Platform,
			URL:           in[i].Url,
		}
	}
	return result
}

func SocialMediaReadModelsToProto(in []*readmodel.SocialMediaReadModel) []*socialmediav1.SocialMediaResponse {
	sm := make([]*socialmediav1.SocialMediaResponse, len(in))
	for i := range in {
		sm[i] = &socialmediav1.SocialMediaResponse{
			Id:       int64(in[i].ID),
			Platform: in[i].Platform,
			Url:      in[i].URL,
		}
	}
	return sm
}
