// Package service
package service

import "github.com/paladignus/actajus/internal/application/dto"

type Token interface {
	GenerateTokenPair(string) (dto.TokenPair, error)
	ValidateAccessToken(string) (dto.TokenClaims, error)
	ValidateRefreshToken(string) (dto.TokenClaims, error)
	RefreshAccessToken(refreshToken string) (dto.TokenPair, error)
}
