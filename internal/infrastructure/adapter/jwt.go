// Package adapter
package adapter

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrBuildToken   = errors.New("failed to build token")
)

type jwtAdapter struct {
	config config.JWTConfig
}

func NewJWTAdapter(config config.JWTConfig) jwtAdapter {
	return jwtAdapter{config: config}
}

func (j jwtAdapter) GenerateTokenPair(IDUser string) (dto.TokenPair, error) {
	accessToken, err := j.generateToken(IDUser, j.config.AccessSecret, j.config.AccessExpire)
	if err != nil {
		return dto.TokenPair{}, err
	}
	refreshToken, err := j.generateToken(IDUser, j.config.RefreshSecret, j.config.RefreshExpire)
	if err != nil {
		return dto.TokenPair{}, err
	}
	return dto.TokenPair{
		AccessToken:  string(accessToken),
		RefreshToken: string(refreshToken),
	}, err
}

func (j jwtAdapter) GenerateResetToken(IDUser string) (dto.TokenRecover, error) {
	token, err := j.generateToken(IDUser, j.config.ResetSecret, j.config.ResetExpire)
	if err != nil {
		return dto.TokenRecover{}, err
	}
	return dto.TokenRecover{ResetToken: string(token)}, nil
}

func (j jwtAdapter) ValidateAccessToken(token string) (dto.TokenClaims, error) {
	return j.validateToken(token, j.config.AccessSecret)
}

func (j jwtAdapter) ValidateRefreshToken(token string) (dto.TokenClaims, error) {
	return j.validateToken(token, j.config.RefreshSecret)
}

func (j jwtAdapter) ValidateResetToken(token string) (dto.TokenClaims, error) {
	return j.validateToken(token, j.config.ResetSecret)
}

func (j jwtAdapter) RefreshAccessToken(token string) (dto.TokenPair, error) {
	claims, err := j.ValidateRefreshToken(token)
	if err != nil {
		return dto.TokenPair{}, err
	}
	return j.GenerateTokenPair(claims.IDUser)
}

func (j jwtAdapter) generateToken(IDUser string, secret string, expire time.Duration) ([]byte, error) {
	now := time.Now()
	token, err := jwt.NewBuilder().
		Subject(IDUser).
		IssuedAt(now).
		Expiration(now.Add(expire)).
		Issuer(j.config.Issuer).
		JwtID(generateJTI()).
		NotBefore(now).
		Build()
	if err != nil {
		return nil, ErrBuildToken
	}
	return jwt.Sign(token, jwt.WithKey(jwa.HS256(), []byte(secret)))
}

func (j jwtAdapter) validateToken(tokenString, secret string) (dto.TokenClaims, error) {
	token, err := jwt.ParseString(tokenString, jwt.WithKey(jwa.HS256(), []byte(secret)))
	if err != nil {
		if errors.Is(err, jwt.TokenExpiredError()) {
			return dto.TokenClaims{}, ErrExpiredToken
		}
		return dto.TokenClaims{}, ErrInvalidToken
	}
	subject, ok := token.Subject()
	if !ok || subject == "" {
		return dto.TokenClaims{}, ErrInvalidToken
	}
	return dto.TokenClaims{
		IDUser: subject,
	}, nil
}

func generateJTI() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
