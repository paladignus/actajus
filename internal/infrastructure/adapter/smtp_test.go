package adapter

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestSMTPEmail(t *testing.T) {
	smtpConfig := config.SMTPConfig{
		Host: "smtp.example.com",
		Port: "587",
		User: "user@example.com",
		Pass: "password",
		From: "from@example.com",
	}
	sut := NewSMTPEmail(smtpConfig)
	assert.Equal(t, "smtp.example.com", sut.config.Host)
	assert.Equal(t, "587", sut.config.Port)
	assert.Equal(t, "user@example.com", sut.config.User)
	emptyConfig := config.SMTPConfig{}
	sut = NewSMTPEmail(emptyConfig)
	assert.Equal(t, "", sut.config.Host)
	assert.Equal(t, "", sut.config.Port)
	ctx := context.Background()
	err := sut.SendEmail(ctx, "to@example.com", "Test Subject", "Test Body")
	assert.Error(t, err)
	err2 := sut.SendEmail(ctx, "", "", "")
	assert.Error(t, err2)
}

