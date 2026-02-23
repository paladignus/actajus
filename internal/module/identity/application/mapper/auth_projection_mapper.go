// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type AuthProjectionMapper struct{}

func NewAuthProjectionMapper() *AuthProjectionMapper {
	return &AuthProjectionMapper{}
}

func (p *AuthProjectionMapper) ProjectTokens(
	idSession identity.IDSession,
	idUser identity.IDUser,
	accessToken string,
	refreshToken string,
	accessExp time.Time,
	refreshExp time.Time,
) *dto.AuthTokensReadModel {
	return &dto.AuthTokensReadModel{
		IDSession:        idSession.Value(),
		IDUser:           idUser.Value(),
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
	}
}
