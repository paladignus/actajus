// Package exception
package exception

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidUserData    = errors.New("invalid user data")
	ErrUserAlreadyExists  = errors.New("user already exists")

	ErrInvalidCPF = errors.New("invalid cpf")
)
