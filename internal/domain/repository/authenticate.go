// Package repository
package repository

type AuthenticateRepository interface {
	Authenticate(username, password string) (string, error)
}
