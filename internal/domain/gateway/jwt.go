// Package gateway
package gateway

import "github.com/paladignus/actajus/internal/application/readmodel"

type JWT interface {
	GenerateTokenPair(string) (readmodel.TokenPairReadModel, error)
	// GenerateResetToken(string) (dto.TokenRecover, error)
	ValidateAccessToken(string) (readmodel.TokenClaimsReadModel, error)
	ValidateRefreshToken(string) (readmodel.TokenClaimsReadModel, error)
	// ValidateResetToken(string) (dto.TokenClaims, error)
	RefreshAccessToken(string) (readmodel.TokenPairReadModel, error)
}
