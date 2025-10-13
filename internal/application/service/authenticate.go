// Package service
package service

type Authenticate interface {
	Authenticate(email string, password string) (string, error)
}
