// Package adapter
package adapter

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/paladignus/actajus/internal/application/readmodel"
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

func (j jwtAdapter) GenerateTokenPair(idUser string) (readmodel.TokenPairReadModel, error) {
	accessToken, err := j.generateToken(idUser, j.config.AccessSecret, j.config.AccessExpire)
	if err != nil {
		return readmodel.TokenPairReadModel{}, err
	}
	refreshToken, err := j.generateToken(idUser, j.config.RefreshSecret, j.config.RefreshExpire)
	if err != nil {
		return readmodel.TokenPairReadModel{}, err
	}
	return readmodel.TokenPairReadModel{
		AccessToken:  string(accessToken),
		RefreshToken: string(refreshToken),
	}, err
}

func (j jwtAdapter) ValidateAccessToken(token string) (readmodel.TokenClaimsReadModel, error) {
	return j.validateToken(token, j.config.AccessSecret)
}

func (j jwtAdapter) ValidateRefreshToken(token string) (readmodel.TokenClaimsReadModel, error) {
	return j.validateToken(token, j.config.RefreshSecret)
}

func (j jwtAdapter) RefreshAccessToken(token string) (readmodel.TokenPairReadModel, error) {
	claims, err := j.ValidateRefreshToken(token)
	if err != nil {
		return readmodel.TokenPairReadModel{}, err
	}
	return j.GenerateTokenPair(claims.IDUser)
}

func (j jwtAdapter) generateToken(idUser string, secret string, expire time.Duration) ([]byte, error) {
	now := time.Now()
	token, err := jwt.NewBuilder().
		Subject(idUser).
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

func (j jwtAdapter) validateToken(tokenString, secret string) (readmodel.TokenClaimsReadModel, error) {
	token, err := jwt.ParseString(tokenString, jwt.WithKey(jwa.HS256(), []byte(secret)))
	if err != nil {
		if errors.Is(err, jwt.TokenExpiredError()) {
			return readmodel.TokenClaimsReadModel{}, ErrExpiredToken
		}
		return readmodel.TokenClaimsReadModel{}, ErrInvalidToken
	}
	subject, ok := token.Subject()
	if !ok || subject == "" {
		return readmodel.TokenClaimsReadModel{}, ErrInvalidToken
	}
	return readmodel.TokenClaimsReadModel{
		IDUser: subject,
	}, nil
}

func generateJTI() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
