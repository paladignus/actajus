// Package spy
package spy

import (
	"github.com/paladignus/actajus/internal/application/dto"
)

// SpyToken implementa gateway.Token
type SpyToken struct {
	Pair       dto.TokenPair
	ResetToken dto.TokenRecover
	Err        error

	CalledWithID string
}

func (s *SpyToken) GenerateTokenPair(idUser string) (dto.TokenPair, error) {
	s.CalledWithID = idUser
	return s.Pair, s.Err
}

func (s *SpyToken) RefreshAccessToken(refreshToken string) (dto.TokenPair, error) {
	return s.Pair, s.Err
}

func (s *SpyToken) ValidateAccessToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

func (s *SpyToken) ValidateRefreshToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

func (s *SpyToken) GenerateResetToken(idUser string) (dto.TokenRecover, error) {
	return s.ResetToken, s.Err
}

func (s *SpyToken) ValidateResetToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}
