// Package service
package service

type IToken interface {
	GenerateToken() (string, error)
}
