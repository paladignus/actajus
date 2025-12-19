// Package spy
package spy

import (
	"github.com/paladignus/actajus/internal/application/dto"
)

// Token implementa gateway.Token
type Token struct {
	Pair dto.TokenPair
	// ResetToken dto.TokenRecover
	Err error

	CalledWithID string
}

func (s *Token) GenerateTokenPair(idUser string) (dto.TokenPair, error) {
	s.CalledWithID = idUser
	return s.Pair, s.Err
}

func (s *Token) RefreshAccessToken(refreshToken string) (dto.TokenPair, error) {
	return s.Pair, s.Err
}

func (s *Token) ValidateAccessToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

func (s *Token) ValidateRefreshToken(tokenString string) (dto.TokenClaims, error) {
	return dto.TokenClaims{}, nil
}

// func (s *Token) GenerateResetToken(idUser string) (dto.TokenRecover, error) {
// 	return s.ResetToken, s.Err
// }
//
// func (s *Token) ValidateResetToken(tokenString string) (dto.TokenClaims, error) {
// 	return dto.TokenClaims{}, nil
// }
