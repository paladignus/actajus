// Package repository
package repository

type Repository interface {
	User() IUser
	Logger() Logger
	Cache() Cache
}
