package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmailRecoverPasswordInput(t *testing.T) {
	input := SendEmailRecoverPasswordInput{
		Address: "john.doe@example.com",
		URL:     "https://example.com/reset?token=abc123",
	}
	t.Run("should return the same address and URL", func(t *testing.T) {
		assert.Equal(t, "john.doe@example.com", input.Address)
		assert.Equal(t, "https://example.com/reset?token=abc123", input.URL)
	})
	t.Run("should return an empty address and URL", func(t *testing.T) {
		input := SendEmailRecoverPasswordInput{}
		assert.Equal(t, "", input.Address)
		assert.Equal(t, "", input.URL)
	})
}

