// Package security
package security

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) BcryptPasswordHasher {
	return BcryptPasswordHasher{
		cost: cost,
	}
}

func (b BcryptPasswordHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(hashedPassword), err
}

func (b BcryptPasswordHasher) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
