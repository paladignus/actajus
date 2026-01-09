package exception

import (
	"testing"
)

func TestErrorVariables(t *testing.T) {
	tests := []struct {
		err      error
		expected string
	}{
		{ErrInvalidFirstName, "first name is invalid"},
		{ErrInvalidLastName, "last name is invalid"},
		{ErrInvalidBirthDate, "birth date is invalid"},
		{ErrInvalidMother, "mother is invalid"},
		{ErrInvalidFather, "father is invalid"},
		{ErrInvalidGender, "gender is invalid"},
		{ErrInvalidMarital, "marital status is invalid"},
		{ErrUserNotFound, "user not found"},
		{ErrInvalidCredentials, "invalid credentials"},
		{ErrInvalidToken, "invalid token"},
		{ErrTokenExpired, "token expired"},
		{ErrUnauthorized, "unauthorized"},
		{ErrForbidden, "forbidden"},
		{ErrInvalidUserData, "invalid user data"},
		{ErrUserAlreadyExists, "user already exists"},
		{ErrInvalidCPF, "invalid cpf"},
		{ErrCPFNotFound, "user not found"},
		{ErrEmailNotFound, "email not found"},
		{ErrEmailAlreadyExists, "email already exists"},
		{ErrEmailNotVerified, "email not verified"},
		{ErrEmailAlreadyVerified, "email already verified"},
		{ErrInvalidEmail, "invalid email"},
		{ErrInvalidPassword, "invalid password"},
		{ErrInvalidRegisteredBy, "registered by is invalid or empty"},
		{ErrInvalidName, "name is invalid"},
		{ErrInvalidTradeName, "trade name is invalid"},
		{ErrInvalidCNPJ, "cnpj is invalid"},
	}
	for _, sut := range tests {
		if sut.err.Error() != sut.expected {
			t.Errorf("Expected error '%s', got '%s'", sut.expected, sut.err.Error())
		}
	}
}
