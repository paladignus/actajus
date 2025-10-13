// Package repository
package repository

type Repository interface {
	Account() Account
	Logger() Logger
	Cache() Cache
}
