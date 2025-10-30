// Package repository
package repository

type Repository interface {
	Authentication() Authentication
	Logger() Logger
	Cache() Cache
}
