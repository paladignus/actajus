// Package dto
package dto

import "time"

type AuthenticationInput struct {
	IDPeople    int        `json:"id_people"`
	Password    string     `json:"password"`
	Avatar      string     `json:"avatar"`
	LastLoginAt time.Time  `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type SignInInput struct {
	CPF      string `json:"cpf"`
	Password string `json:"password"`
}

type SignInOutput struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	IDUser       string   `json:"id_user"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
}

type Role struct {
	IDRole int    `json:"-"`
	Name   string `json:"name"`
}

type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type GetEmailByCPFInput struct {
	CPF string `json:"cpf"`
}

type GetEmailByCPFOutput struct {
	Email string `json:"email"`
}

type RecoverPasswordInput struct {
	Email string `json:"email"`
}
