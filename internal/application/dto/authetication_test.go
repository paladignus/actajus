package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignInInput(t *testing.T) {
	input := SignInInput{
		CPF:      "12345678901",
		Password: "password123",
	}
	t.Run("should return the same CPF and password", func(t *testing.T) {
		assert.Equal(t, "12345678901", input.CPF)
		assert.Equal(t, "password123", input.Password)
	})
	t.Run("should return an empty CPF and password", func(t *testing.T) {
		input := SignInInput{}
		assert.Empty(t, input.CPF)
		assert.Empty(t, input.Password)
	})
}

func TestSignInOutput(t *testing.T) {
	output := SignInOutput{
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-456",
		IDUser:       "user-id-789",
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john.doe@example.com",
		Roles:        []string{"admin", "user"},
	}
	t.Run("should return the same access token, refresh token, ID user, first name, last name, email and roles", func(t *testing.T) {
		assert.Equal(t, "access-token-123", output.AccessToken)
		assert.Equal(t, "refresh-token-456", output.RefreshToken)
		assert.Equal(t, "user-id-789", output.IDUser)
		assert.Equal(t, "John", output.FirstName)
		assert.Equal(t, "Doe", output.LastName)
		assert.Equal(t, "john.doe@example.com", output.Email)
		expectedRoles := []string{"admin", "user"}
		assert.Equal(t, expectedRoles, output.Roles)
	})
	t.Run("should return an empty access token, refresh token, ID user, first name, last name, email and roles", func(t *testing.T) {
		output := SignInOutput{}
		assert.Empty(t, output.AccessToken)
		assert.Empty(t, output.RefreshToken)
		assert.Empty(t, output.IDUser)
		assert.Empty(t, output.FirstName)
		assert.Empty(t, output.LastName)
		assert.Empty(t, output.Email)
		assert.Empty(t, output.Roles)
	})
}

func TestRole(t *testing.T) {
	role := Role{
		IDRole: 1,
		Name:   "admin",
	}
	t.Run("should return the same ID role and name", func(t *testing.T) {
		assert.Equal(t, 1, role.IDRole)
		assert.Equal(t, "admin", role.Name)
	})
	t.Run("should return an empty ID role and name", func(t *testing.T) {
		role := Role{}
		assert.Equal(t, 0, role.IDRole)
		assert.Equal(t, "", role.Name)
	})
}

func TestPermission(t *testing.T) {
	permission := Permission{
		Resource: "user",
		Action:   "read",
	}
	t.Run("should return the same resource and action", func(t *testing.T) {
		assert.Equal(t, "user", permission.Resource)
		assert.Equal(t, "read", permission.Action)
	})
	t.Run("should return an empty resource and action", func(t *testing.T) {
		permission := Permission{}
		assert.Equal(t, "", permission.Resource)
		assert.Equal(t, "", permission.Action)
	})
}

func TestGetEmailByCPFInput(t *testing.T) {
	input := GetEmailByCPFInput{
		CPF: "12345678901",
	}
	t.Run("should return the same CPF", func(t *testing.T) {
		assert.Equal(t, "12345678901", input.CPF)
	})
	t.Run("should return an empty CPF", func(t *testing.T) {
		input := GetEmailByCPFInput{}
		assert.Equal(t, "", input.CPF)
	})
}

func TestGetEmailByCPFOutput(t *testing.T) {
	output := GetEmailByCPFOutput{
		Email: "john.doe@example.com",
	}
	t.Run("should return the same email", func(t *testing.T) {
		assert.Equal(t, "john.doe@example.com", output.Email)
	})
	t.Run("should return an empty email", func(t *testing.T) {
		output := GetEmailByCPFOutput{}
		assert.Equal(t, "", output.Email)
	})
}

func TestRecoverPasswordInput(t *testing.T) {
	input := RecoverPasswordInput{
		Email: "john.doe@example.com",
	}
	t.Run("should return the same email", func(t *testing.T) {
		assert.Equal(t, "john.doe@example.com", input.Email)
	})
	t.Run("should return an empty email", func(t *testing.T) {
		input := RecoverPasswordInput{}
		assert.Equal(t, "", input.Email)
	})
}

