// Package spy
package spy

import (
	"errors"
	"time"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token has expired")
	ErrWrongTokenType = errors.New("wrong token type")
)

type JWTAdapter struct {
	SecretKey        string
	ExpiresIn        time.Duration
	CallsCount       int
	ErrGenerateToken error
}

func NewJWTAdapter(secretKey string, expiresIn time.Duration) JWTAdapter {
	return JWTAdapter{SecretKey: secretKey, ExpiresIn: expiresIn}
}

func (j *JWTAdapter) NewToken() (string, error) {
	j.CallsCount++
	return "valid_token", j.ErrGenerateToken
}
