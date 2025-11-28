package messaging

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/repository"
)

// MockSMTP is a mock implementation of gateway.SMTP for testing
type MockSMTP struct {
	sendEmailError error
}

func (m *MockSMTP) SendEmail(ctx context.Context, to, subject, body string) error {
	return m.sendEmailError
}

// MockLogger is a mock implementation of repository.Logger for testing
type MockLogger struct {
	logCalls []string
}

func (m *MockLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	m.logCalls = append(m.logCalls, "debug")
}

func (m *MockLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	m.logCalls = append(m.logCalls, "info")
}

func (m *MockLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	m.logCalls = append(m.logCalls, "warn")
}

func (m *MockLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	m.logCalls = append(m.logCalls, "error")
}

func (m *MockLogger) With(args ...interface{}) repository.Logger {
	return m
}

func (m *MockLogger) WithError(err error) repository.Logger {
	return m
}

// TestRecoverPassword tests the RecoverPassword messaging functionality
func TestRecoverPassword(t *testing.T) {
	ctx := context.Background()

	t.Run("successful email sending", func(t *testing.T) {
		mockLogger := &MockLogger{}
		mockSMTP := &MockSMTP{}

		recoverPassword := NewRecoverPassword(mockLogger, mockSMTP)

		// Create a PasswordResetRequestedEvent
		testEvent := event.NewPasswordResetRequestedEvent(
			"user-123",
			"test@example.com",
			"https://example.com/reset",
		)

		err := recoverPassword.Handle(ctx, &testEvent)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check if CanHandle works correctly
		if !recoverPassword.CanHandle(&testEvent) {
			t.Error("Expected CanHandle to return true for PasswordResetRequestedEvent")
		}
	})

	t.Run("wrong event type", func(t *testing.T) {
		mockLogger := &MockLogger{}
		mockSMTP := &MockSMTP{}

		recoverPassword := NewRecoverPassword(mockLogger, mockSMTP)

		// Create a custom event with a different name - removing this problematic approach
		// Instead, let's just test the PasswordResetRequestedEvent check directly

		// Create a proper event that doesn't match the expected type
		differentEvent := event.NewEvent("different.event", "agg-id", "1.0")

		err := recoverPassword.Handle(ctx, differentEvent)

		// This should fail because it's not a *event.PasswordResetRequestedEvent
		if err == nil {
			t.Error("Expected error for wrong event type, got nil")
		}
	})

	t.Run("email sending fails", func(t *testing.T) {
		mockLogger := &MockLogger{}
		mockSMTP := &MockSMTP{
			sendEmailError: fmt.Errorf("SMTP error for testing"), // Using a generic error for testing
		}

		recoverPassword := NewRecoverPassword(mockLogger, mockSMTP)

		// Create a PasswordResetRequestedEvent
		testEvent := event.NewPasswordResetRequestedEvent(
			"user-123",
			"test@example.com",
			"https://example.com/reset",
		)

		err := recoverPassword.Handle(ctx, &testEvent)

		if err == nil {
			t.Error("Expected error when email sending fails, got nil")
		}
	})

	t.Run("CanHandle with correct event name", func(t *testing.T) {
		mockLogger := &MockLogger{}
		mockSMTP := &MockSMTP{}

		recoverPassword := NewRecoverPassword(mockLogger, mockSMTP)

		// Create a PasswordResetRequestedEvent
		testEvent := event.NewPasswordResetRequestedEvent(
			"user-123",
			"test@example.com",
			"https://example.com/reset",
		)

		if !recoverPassword.CanHandle(&testEvent) {
			t.Error("Expected CanHandle to return true for PasswordResetRequestedEvent")
		}
	})

	t.Run("CanHandle with different event name", func(t *testing.T) {
		mockLogger := &MockLogger{}
		mockSMTP := &MockSMTP{}

		recoverPassword := NewRecoverPassword(mockLogger, mockSMTP)

		// Create a mock event with different name using the Event struct
		differentEvent := event.NewEvent("different.event", "agg-id", "1.0")

		if recoverPassword.CanHandle(differentEvent) {
			t.Error("Expected CanHandle to return false for different event name")
		}
	})

	// Test body function
	t.Run("body function generates correct HTML", func(t *testing.T) {
		resetURL := "https://example.com/reset"
		result := body(resetURL)

		if result == "" {
			t.Error("Expected body to not be empty")
		}

		// Check if the URL is properly included in the body
		if len(result) < 100 { // Basic check for reasonable length
			t.Error("Expected body to have reasonable length")
		}

		if !strings.Contains(result, resetURL) {
			t.Error("Expected body to contain the reset URL")
		}

		if !strings.Contains(result, "Recuperar Senha") {
			t.Error("Expected body to contain password recovery text")
		}
	})
}

