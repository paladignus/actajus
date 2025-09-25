// Package adapter
package adapter

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAdapter struct{}

func NewJWTAdapter() JWTAdapter {
	return JWTAdapter{}
}

func (j *JWTAdapter) NewToken(ctx context.Context, secreyKey string, expiresIn time.Duration) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "your-issuer",
		"exp": time.Now().Add(expiresIn).Unix(),
		"iat": time.Now().Unix(),
		"sub": "user-id",
	})

	return token.SignedString([]byte(secreyKey))
}
