// Package service
package service

type ITokenGenerator interface {
	Generate() (string, error)
}
