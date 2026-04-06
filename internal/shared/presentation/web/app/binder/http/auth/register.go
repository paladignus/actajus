package authbinder

import (
	"net/http"
	"strconv"
	"strings"
)

type RegisterForm struct {
	FirstName string
	LastName  string
	Birthday  string
	GenderID  uint
	Email     string
	Password  string
}

type VerifyEmailQuery struct {
	IDVerification int64
	Token          string
}

func BindRegisterForm(r *http.Request) (RegisterForm, error) {
	if err := r.ParseForm(); err != nil {
		return RegisterForm{}, err
	}
	genderID, _ := strconv.ParseUint(strings.TrimSpace(r.FormValue("gender_id")), 10, 32)
	return RegisterForm{
		FirstName: strings.TrimSpace(r.FormValue("first_name")),
		LastName:  strings.TrimSpace(r.FormValue("last_name")),
		Birthday:  strings.TrimSpace(r.FormValue("birthday")),
		GenderID:  uint(genderID),
		Email:     strings.TrimSpace(r.FormValue("email")),
		Password:  r.FormValue("password"),
	}, nil
}

func BindVerifyEmailQuery(r *http.Request) VerifyEmailQuery {
	idVerification, _ := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("id_verification")), 10, 64)
	return VerifyEmailQuery{
		IDVerification: idVerification,
		Token:          strings.TrimSpace(r.URL.Query().Get("token")),
	}
}
