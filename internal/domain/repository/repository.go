// Package repository
package repository

type Repository interface {
	Account() Account
	Token() Token
	Logger() Logger
	Cache() Cache
}
