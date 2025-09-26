// Package adapter
package adapter

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtAdapter struct {
	accessSecre   []byte
	refreshSecret []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

type customClaims struct {
	UserID string
	jwt.RegisteredClaims
}

func NewJWTAdapter() jwtAdapter {
	return jwtAdapter{}
}

func (j *jwtAdapter) generateToken(userID string, secret []byte, expiresIn time.Duration) (string, error) {
	jti := uuid.New().String()
	claims := customClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *jwtAdapter) validateToken(tokenString string, secret []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inválido")
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		//cache token, ver depois
		return claims.UserID, nil
	}
	return "", errors.New("token inválido")
}

func (j *jwtAdapter) revokeToken(tokenString string, secret []byte) error {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(*customClaims); ok {
		// delete token cache posterior
		return nil
	}
	return errors.New("token inválido")
}

func (j *jwtAdapter) GenereteAccessToken(userID string) (string, error) {
	return j.generateToken(userID, j.accessSecre, j.accessExpiry)
}

func (j *jwtAdapter) GenereteRefreshToken(userID string) (string, error) {
	return j.generateToken(userID, j.refreshSecret, j.refreshExpiry)
}

func (j *jwtAdapter) ValidateAccessToken(token string) (string, error) {
	return j.validateToken(token, j.accessSecre)
}

func (j *jwtAdapter) ValidateRefreshToken(token string) (string, error) {
	return j.validateToken(token, j.refreshSecret)
}

func (j *jwtAdapter) RevokeAccessToken(token string) error {
	return j.revokeToken(token, j.accessSecre)
}

func (j *jwtAdapter) RevokeRefreshToken(token string) error {
	return j.revokeToken(token, j.refreshSecret)
}
