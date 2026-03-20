// Package spy
package spy

import (
	"github.com/paladignus/actajus/internal/application/readmodel"
)

type JWT struct {
	Pair readmodel.TokenPairReadModel
	// ResetToken dto.TokenRecover
	Err error

	CalledWithID string
}

func (s *JWT) GenerateTokenPair(idUser string) (readmodel.TokenPairReadModel, error) {
	s.CalledWithID = idUser
	return s.Pair, s.Err
}

func (s *JWT) RefreshAccessToken(refreshToken string) (readmodel.TokenPairReadModel, error) {
	return s.Pair, s.Err
}

func (s *JWT) ValidateAccessToken(tokenString string) (readmodel.TokenClaimsReadModel, error) {
	return readmodel.TokenClaimsReadModel{}, nil
}

func (s *JWT) ValidateRefreshToken(tokenString string) (readmodel.TokenClaimsReadModel, error) {
	return readmodel.TokenClaimsReadModel{}, nil
}

// func (s *Token) GenerateResetToken(idUser string) (dto.TokenRecover, error) {
// 	return s.ResetToken, s.Err
// }
//
// func (s *Token) ValidateResetToken(tokenString string) (dto.TokenClaims, error) {
// 	return readmodel.TokenClaimsReadModel{}, nil
// }
