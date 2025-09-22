// Package usecase
package usecase

import "github.com/paladignus/actajus/internal/domain/repository"

type SignIn struct {
	repository repository.AuthenticateRepository
}

func NewSignIn(repository repository.AuthenticateRepository) SignIn {
	return SignIn{repository}
}

func (s SignIn) Execute(username, password string) (string, error) {
	return s.repository.Authenticate(username, password)
}
