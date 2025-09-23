// Package repository
package repository

type AuthenticateRepository interface {
	Authenticate(cpf, password string) (string, error)
}
