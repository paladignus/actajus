// Package service
package service

type Password interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}
