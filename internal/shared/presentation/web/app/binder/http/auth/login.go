package authbinder

import (
	"net/http"
	"strings"
)

type LoginForm struct {
	Email    string
	Password string
}

func BindLoginForm(r *http.Request) (LoginForm, error) {
	if err := r.ParseForm(); err != nil {
		return LoginForm{}, err
	}
	return LoginForm{
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: r.FormValue("password"),
	}, nil
}
