// Package spy
package spy

import "errors"

type AuthenticateSpyPersistence struct {
	Username   string
	Password   string
	CallsCount int
}

func (a *AuthenticateSpyPersistence) Authenticate(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password cannot be empty")
	}
	a.Username = username
	a.Password = password
	a.CallsCount++
	return "user_id", nil
}
