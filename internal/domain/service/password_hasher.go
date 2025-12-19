// Package service
package service

type PasswordHasher interface {
	Hash(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}
