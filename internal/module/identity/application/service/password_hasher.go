// Package service
package service

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(encodedHash string, plain string) error
	NeedsRehash(encodedHash string) bool
}
