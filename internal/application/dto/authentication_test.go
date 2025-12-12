package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignInInput(t *testing.T) {
	sut := SignInInput{
		CPF:      "12345678901",
		Password: "password123",
	}
	t.Run("should return the same CPF and password", func(t *testing.T) {
		assert.Equal(t, "12345678901", sut.CPF)
		assert.Equal(t, "password123", sut.Password)
	})
	t.Run("should return an empty CPF and password", func(t *testing.T) {
		sut := SignInInput{}
		assert.Empty(t, sut.CPF)
		assert.Empty(t, sut.Password)
	})
}

func TestSignInOutput(t *testing.T) {
	sut := SignInOutput{
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-456",
		IDUser:       "user-id-789",
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john.doe@example.com",
		Roles:        []string{"admin", "user"},
	}
	t.Run("should return the same access token, refresh token, ID user, first name, last name, email and roles", func(t *testing.T) {
		assert.Equal(t, "access-token-123", sut.AccessToken)
		assert.Equal(t, "refresh-token-456", sut.RefreshToken)
		assert.Equal(t, "user-id-789", sut.IDUser)
		assert.Equal(t, "John", sut.FirstName)
		assert.Equal(t, "Doe", sut.LastName)
		assert.Equal(t, "john.doe@example.com", sut.Email)
		expectedRoles := []string{"admin", "user"}
		assert.Equal(t, expectedRoles, sut.Roles)
	})
	t.Run("should return an empty access token, refresh token, ID user, first name, last name, email and roles", func(t *testing.T) {
		sut := SignInOutput{}
		assert.Empty(t, sut.AccessToken)
		assert.Empty(t, sut.RefreshToken)
		assert.Empty(t, sut.IDUser)
		assert.Empty(t, sut.FirstName)
		assert.Empty(t, sut.LastName)
		assert.Empty(t, sut.Email)
		assert.Empty(t, sut.Roles)
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
	sut := GetEmailByCPFInput{
		CPF: "12345678901",
	}
	t.Run("should return the same CPF", func(t *testing.T) {
		assert.Equal(t, "12345678901", sut.CPF)
	})
	t.Run("should return an empty CPF", func(t *testing.T) {
		sut := GetEmailByCPFInput{}
		assert.Equal(t, "", sut.CPF)
	})
}

func TestGetEmailByCPFOutput(t *testing.T) {
	sut := GetEmailByCPFOutput{
		Email: "john.doe@example.com",
	}
	t.Run("should return the same email", func(t *testing.T) {
		assert.Equal(t, "john.doe@example.com", sut.Email)
	})
	t.Run("should return an empty email", func(t *testing.T) {
		sut := GetEmailByCPFOutput{}
		assert.Equal(t, "", sut.Email)
	})
}

func TestRecoverPasswordInput(t *testing.T) {
	sut := RecoverPasswordInput{
		Email: "john.doe@example.com",
	}
	t.Run("should return the same email", func(t *testing.T) {
		assert.Equal(t, "john.doe@example.com", sut.Email)
	})
	t.Run("should return an empty email", func(t *testing.T) {
		sut := RecoverPasswordInput{}
		assert.Equal(t, "", sut.Email)
	})
}
