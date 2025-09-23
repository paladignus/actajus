// Package spy
package spy

import "errors"

type AuthenticateSpyPersistence struct {
	CPF        string
	Password   string
	CallsCount int
}

func (a *AuthenticateSpyPersistence) Authenticate(cpf, password string) (string, error) {
	if cpf == "" || password == "" {
		return "", errors.New("cpf and password cannot be empty")
	}
	a.CPF = cpf
	a.Password = password
	a.CallsCount++
	return "user_id", nil
}
