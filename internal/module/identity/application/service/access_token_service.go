// Package service
package service

import "github.com/paladignus/actajus/internal/module/identity/application/readmodel"

type AccessTokenService interface {
	Sign(claims readmodel.AccessTokenClaims) (string, error)
	Verify(token string) (readmodel.AccessTokenClaims, error)
}
