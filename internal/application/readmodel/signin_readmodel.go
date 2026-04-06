// Package readmodel
package readmodel

type SignInReadModel struct {
	IDUser       int
	AccessToken  string
	RefreshToken string
	FirstName    string
	LastName     string
	Email        string
	Roles        []string
}
