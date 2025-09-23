// Package dto
package dto

type AuthenticatedInput struct {
	CPF      string
	Password string
}

type AuthenticatedOutput struct {
	AccessToken  string
	RefreshToken string
	ID           string
	FirstName    string
	LastName     string
	Email        string
	Roles        []string
}
