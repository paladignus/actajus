// Package dto
package dto

type AuthenticatedInput struct {
	CPF      string
	Password string
}

type Role struct {
	IDRole int    `json:"-"`
	Name   string `json:"name"`
}

type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type AuthenticatedOutput struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	IDPerson     string   `json:"id_person"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
}
