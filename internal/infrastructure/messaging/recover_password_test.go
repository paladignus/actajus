package messaging

import (
	"context"
	"fmt"
	"testing"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSMTP struct {
	sendEmailError error
}

func (m *MockSMTP) SendEmail(ctx context.Context, to, subject, body string) error {
	return m.sendEmailError
}

func TestRecoverPassword(t *testing.T) {
	ctx := context.Background()
	logger := &spy.Logger{}
	mockSMTP := &MockSMTP{}
	sut := NewRecoverPassword(logger, mockSMTP)
	t.Run("should success email sending", func(t *testing.T) {
		testEvent := event.NewPasswordResetRequestedEvent(
			"user-123",
			"test@example.com",
			"https://example.com/reset",
		)
		logger.On("Info", ctx, "email sent successfully", "to", "test@example.com")
		err := sut.Handle(ctx, &testEvent)
		assert.NoError(t, err)
		assert.True(t, sut.CanHandle(&testEvent))
	})

	t.Run("should wrong event type", func(t *testing.T) {
		differentEvent := event.NewEvent("different.event", "agg-id", "1.0")
		logger.On("Error", ctx, "invalid event type", "expected", "*PasswordResetRequestedEvent", "got", mock.Anything)
		err := sut.Handle(ctx, differentEvent)
		assert.Error(t, err)
	})

	t.Run("should email sending fails", func(t *testing.T) {
		mockSMTP := &MockSMTP{
			sendEmailError: fmt.Errorf("SMTP error for testing"),
		}
		logger.On("Error", ctx, "failed to send email", "to", "test@example.com", "error", mock.Anything)
		sut := NewRecoverPassword(logger, mockSMTP)
		testEvent := event.NewPasswordResetRequestedEvent(
			"user-123",
			"test@example.com",
			"https://example.com/reset",
		)
		err := sut.Handle(ctx, &testEvent)
		assert.Error(t, err)
	})

	t.Run("should CanHandle with correct event name", func(t *testing.T) {
		testEvent := event.NewPasswordResetRequestedEvent(
			"user-123",
			"test@example.com",
			"https://example.com/reset",
		)
		assert.True(t, sut.CanHandle(&testEvent))
	})

	t.Run("should CanHandle with different event name", func(t *testing.T) {
		differentEvent := event.NewEvent("different.event", "agg-id", "1.0")
		assert.False(t, sut.CanHandle(differentEvent))
	})

	t.Run("should body function generates correct HTML", func(t *testing.T) {
		resetURL := "https://example.com/reset"
		result := body(resetURL)
		assert.NotEmpty(t, result)
		assert.GreaterOrEqual(t, len(result), 100)
		assert.Contains(t, result, resetURL)
		assert.Contains(t, result, "Recuperar Senha")
	})
}
