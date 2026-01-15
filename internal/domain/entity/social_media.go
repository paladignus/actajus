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

func NewSocialMedia(input dto.SocialMedia) (SocialMedia, error) {
	s := SocialMedia{
		IDEnterprise: input.IDEnterprise,
		Name:         vo.Text(input.Name),
		URL:          input.URL,
	}
	if !s.Name.IsValid() {
		return s, exception.ErrInvalidName
	}
	return s, nil
}
