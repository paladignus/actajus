// Package domain
package domain

type UserPasswordHashRepository interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
