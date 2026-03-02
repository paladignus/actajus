// Package security
package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
)

type RefreshTokenService struct {
	nBytes int
}

func NewRefreshTokenService() *RefreshTokenService {
	return &RefreshTokenService{nBytes: 32}
}

func NewRefreshTokenServiceWithBytes(n int) (*RefreshTokenService, error) {
	if n < 32 {
		return nil, errors.New("refresh token bytes must be >= 32")
	}
	return &RefreshTokenService{nBytes: n}, nil
}

func (s *RefreshTokenService) Generate() (token string, hash [32]byte, err error) {
	b := make([]byte, s.nBytes)
	if _, err = rand.Read(b); err != nil {
		return "", [32]byte{}, err
	}
	token = base64.RawURLEncoding.EncodeToString(b)
	hash = sha256.Sum256(b)
	return token, hash, nil
}

func (s *RefreshTokenService) Compare(token string, expectedHash [32]byte) bool {
	b, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return false
	}
	got := sha256.Sum256(b)
	return subtle.ConstantTimeCompare(got[:], expectedHash[:]) == 1
}

func (s *RefreshTokenService) Hash(token string) ([32]byte, bool) {
	b, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return [32]byte{}, false
	}
	h := sha256.Sum256(b)
	return h, true
}
