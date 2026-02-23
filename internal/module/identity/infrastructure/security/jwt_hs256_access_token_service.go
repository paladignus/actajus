// Package security
package security

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/paladignus/actajus/internal/module/identity/application/dto"
)

type HS256AccessTokenService struct {
	secret   []byte
	issuer   string
	audience string
}

func NewHS256AccessTokenService(secret string, issuer string, audience string) (*HS256AccessTokenService, error) {
	if secret == "" {
		return nil, errors.New("jwt access secret is empty")
	}
	if issuer == "" {
		return nil, errors.New("jwt issuer is empty")
	}
	if audience == "" {
		return nil, errors.New("jwt audience is empty")
	}
	return &HS256AccessTokenService{
		secret:   []byte(secret),
		issuer:   issuer,
		audience: audience,
	}, nil
}

func (s *HS256AccessTokenService) Sign(claims dto.AccessTokenClaims) (string, error) {
	// now := time.Now()
	// "iat": now.Unix(),
	mapClaims := jwt.MapClaims{
		"iss": s.issuer,
		"aud": s.audience,
		"iat": claims.IssuedAt.Unix(),
		"exp": claims.ExpiresAt.Unix(),
		"sub": strconv.FormatInt(claims.IDUser, 10),
		"sid": strconv.FormatInt(claims.IDSession, 10),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tok.Header["typ"] = "JWT"
	signed, err := tok.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *HS256AccessTokenService) Verify(token string) (dto.AccessTokenClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), // trava HS256
		jwt.WithIssuer(s.issuer),
		jwt.WithAudience(s.audience),
	)
	claims := jwt.MapClaims{}
	_, err := parser.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return s.secret, nil
	})
	if err != nil {
		return dto.AccessTokenClaims{}, err
	}
	sub, _ := claims["sub"].(string)
	sid, _ := claims["sid"].(string)
	if sub == "" || sid == "" {
		return dto.AccessTokenClaims{}, errors.New("missing sub or sid")
	}
	uid64, err := strconv.ParseInt(sub, 10, 64)
	if err != nil || uid64 <= 0 {
		return dto.AccessTokenClaims{}, errors.New("invalid sub")
	}
	sid64, err := strconv.ParseInt(sid, 10, 64)
	if err != nil || sid64 <= 0 {
		return dto.AccessTokenClaims{}, errors.New("invalid sid")
	}
	iat, err := claims.GetIssuedAt()
	if err != nil || iat == nil {
		return dto.AccessTokenClaims{}, errors.New("invalid iat")
	}
	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return dto.AccessTokenClaims{}, errors.New("invalid exp")
	}
	return dto.AccessTokenClaims{
		IDUser:    uid64,
		IDSession: sid64,
		Issuer:    s.issuer,
		Audience:  s.audience,
		IssuedAt:  iat.Time,
		ExpiresAt: exp.Time,
	}, nil
}
