// Package gateway
package gateway

import "github.com/paladignus/actajus/internal/application/dto"

type Token interface {
	GenerateTokenPair(string) (dto.TokenPair, error)
	GenerateResetToken(string) (dto.TokenRecover, error)
	ValidateAccessToken(string) (dto.TokenClaims, error)
	ValidateRefreshToken(string) (dto.TokenClaims, error)
	ValidateResetToken(string) (dto.TokenClaims, error)
	RefreshAccessToken(string) (dto.TokenPair, error)
}
