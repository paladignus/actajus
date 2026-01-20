package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmailRecoverPasswordInput(t *testing.T) {
	sut := SendEmailRecoverPasswordInput{
		Address: "john.doe@example.com",
		URL:     "https://example.com/reset?token=abc123",
	}
	t.Run("should return the same address and URL", func(t *testing.T) {
		assert.Equal(t, "john.doe@example.com", sut.Address)
		assert.Equal(t, "https://example.com/reset?token=abc123", sut.URL)
	})
	t.Run("should return an empty address and URL", func(t *testing.T) {
		sut := SendEmailRecoverPasswordInput{}
		assert.Equal(t, "", sut.Address)
		assert.Equal(t, "", sut.URL)
	})
}
