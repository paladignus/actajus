// Package exception
package exception

import "errors"

var (
	ErrInvalidFirstName = errors.New("first name is invalid")
	ErrInvalidLastName  = errors.New("last name is invalid")
	ErrInvalidBirthDate = errors.New("birth date is invalid")
	ErrInvalidMother    = errors.New("mother is invalid")
	ErrInvalidFather    = errors.New("father is invalid")
	ErrInvalidGender    = errors.New("gender is invalid")
	ErrInvalidMarital   = errors.New("marital status is invalid")

	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidUserData    = errors.New("invalid user data")
	ErrUserAlreadyExists  = errors.New("user already exists")

	ErrInvalidCPF  = errors.New("invalid cpf")
	ErrCPFNotFound = errors.New("user not found")

	ErrEmailNotFound        = errors.New("email not found")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrEmailNotVerified     = errors.New("email not verified")
	ErrEmailAlreadyVerified = errors.New("email already verified")
	ErrInvalidEmail         = errors.New("invalid email")

	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidAvatar      = errors.New("invalid avatar")
	ErrInvalidLastLoginAt = errors.New("invalid last login at")

	ErrInvalidRegisteredBy = errors.New("registered by is invalid or empty")
	ErrInvalidName         = errors.New("name is invalid")
	ErrInvalidTradeName    = errors.New("trade name is invalid")
	ErrInvalidCNPJ         = errors.New("cnpj is invalid")

	ErrInvalidIDUser       = errors.New("id user is invalid")
	ErrInvalidZip          = errors.New("zip is invalid")
	ErrInvalidTitle        = errors.New("title is invalid")
	ErrInvalidStreet       = errors.New("street is invalid")
	ErrInvalidNeighborhood = errors.New("neighborhood is invalid")
	ErrInvalidCity         = errors.New("city is invalid")
	ErrInvalidState        = errors.New("state is invalid")
	ErrInvalidCountry      = errors.New("country is invalid")
)
