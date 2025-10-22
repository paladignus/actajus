// Package adapter
package adapter

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrRevokedToken     = errors.New("token has been revoked")
	ErrInvalidTokenType = errors.New("invalid token type")
)

type customClaims struct {
	IDUser    string
	TokenType string
	jwt.RegisteredClaims
}

type jwtAdapter struct {
	config config.JWTConfig
}

func NewJWTAdapter(config config.JWTConfig) gateway.Token {
	return jwtAdapter{config}
}

func (j jwtAdapter) GenerateTokenPair(IDUser string) (dto.TokenPair, error) {
	accessToken, err := j.generateToken(
		IDUser,
		"access",
		j.config.AccessSecret,
		j.config.AccessExpiry,
	)
	if err != nil {
		return dto.TokenPair{}, err
	}
	refreshToken, err := j.generateToken(
		IDUser,
		"refresh",
		j.config.RefreshSecret,
		j.config.RefreshExpiry,
	)
	if err != nil {
		return dto.TokenPair{}, err
	}
	return dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j jwtAdapter) generateToken(IDUser, tokenType, secret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	jti := j.generateJTI()
	claims := customClaims{
		IDUser:    IDUser,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (j jwtAdapter) ValidateAccessToken(tokenString string) (dto.TokenClaims, error) {
	return j.validateToken(tokenString, "access", j.config.AccessSecret)
}

func (j jwtAdapter) ValidateRefreshToken(tokenString string) (dto.TokenClaims, error) {
	return j.validateToken(tokenString, "refresh", j.config.RefreshSecret)
}

func (j jwtAdapter) validateToken(tokenString, expectedType, secret string) (dto.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return dto.TokenClaims{}, ErrExpiredToken
		}
		return dto.TokenClaims{}, ErrInvalidToken
	}
	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return dto.TokenClaims{}, ErrInvalidToken
	}
	if claims.TokenType != expectedType {
		return dto.TokenClaims{}, ErrInvalidTokenType
	}
	return dto.TokenClaims{
		IDUser:    claims.IDUser,
		TokenType: claims.TokenType,
	}, nil
}

func (j jwtAdapter) RefreshAccessToken(refreshToken string) (dto.TokenPair, error) {
	claims, err := j.ValidateRefreshToken(refreshToken)
	if err != nil {
		return dto.TokenPair{}, err
	}
	return j.GenerateTokenPair(claims.IDUser)
}

func (j jwtAdapter) generateJTI() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
