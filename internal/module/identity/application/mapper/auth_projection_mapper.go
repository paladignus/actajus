// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
)

type AuthProjectionMapper struct{}

func NewAuthProjectionMapper() *AuthProjectionMapper {
	return &AuthProjectionMapper{}
}

func (p *AuthProjectionMapper) ProjectTokens(
	sid int64,
	uid int64,
	accessToken string,
	refreshToken string,
	accessExp time.Time,
	refreshExp time.Time,
) *readmodel.AuthTokensReadModel {
	return &readmodel.AuthTokensReadModel{
		IDSession:        sid,
		IDUser:           uid,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
	}
}
