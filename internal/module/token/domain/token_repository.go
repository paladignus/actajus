// Package domain
package domain

type TokenRepository interface {
	// Comentários para lembra os nomes que IA deu
	// SignAccessToken(userID, sessionID string, role []string) (string, time.Time, error)
	// VerifyAccessToken(token string) (claims []string, err error)
	GenerateToken(claims AccessClaims) (token string, err error)
	VerifyToken(token string) (AccessClaims, error)
}
