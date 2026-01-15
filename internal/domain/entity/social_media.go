// Package entity
package entity

import (
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type SocialMedia struct {
	IDSocialMedia uint
	IDEnterprise  uint
	Name          vo.Text
	URL           string
}

func NewSocialMedia(input dto.SocialMedia) SocialMedia {
	return SocialMedia{
		IDSocialMedia: input.IDSocialMedia,
		IDEnterprise:  input.IDEnterprise,
		Name:          vo.Text(input.Name),
		URL:           input.URL,
	}
}

func (s SocialMedia) Create() error {
	return s.validate()
}

func (s SocialMedia) Update() error {
	if s.IDSocialMedia == 0 {
		return exception.ErrInvalidIDSocialMedia
	}
	return s.validate()
}

func (s SocialMedia) validate() error {
	if !s.Name.IsValid() {
		return exception.ErrInvalidName
	}
	return nil
}
