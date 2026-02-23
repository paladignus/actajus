// Package service
package service

type RefreshTokenService interface {
	Generate() (token string, hash [32]byte, err error)
	Compare(token string, expectedHash [32]byte) bool
}
