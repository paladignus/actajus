// Package spy
package spy

import (
	"github.com/paladignus/actajus/internal/application/dto"
)

type JWT struct {
	Pair dto.TokenPair
	// ResetToken dto.TokenRecover
	Err error

	CalledWithID string
}

func (s *JWT) GenerateTokenPair(idUser string) (dto.TokenPair, error) {
	s.CalledWithID = idUser
	return s.Pair, s.Err
}

func (s *JWT) RefreshAccessToken(refreshToken string) (dto.TokenPair, error) {
	return s.Pair, s.Err
}

func (s *JWT) ValidateAccessToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

func (s *JWT) ValidateRefreshToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

// func (s *Token) GenerateResetToken(idUser string) (dto.TokenRecover, error) {
// 	return s.ResetToken, s.Err
// }
//
// func (s *Token) ValidateResetToken(tokenString string) (dto.TokenClaims, error) {
// 	return dto.TokenClaims{}, nil
// }
