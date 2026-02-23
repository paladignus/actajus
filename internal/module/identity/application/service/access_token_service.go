// Package service
package service

import "github.com/paladignus/actajus/internal/module/identity/application/dto"

type AccessTokenService interface {
	Sign(claims dto.AccessTokenClaims) (string, error)
	Verify(token string) (dto.AccessTokenClaims, error)
}
