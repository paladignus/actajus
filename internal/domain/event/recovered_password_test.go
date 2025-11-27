// Package event
package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasswordResetRequestedEvent(t *testing.T) {
	event := NewPasswordResetRequestedEvent(
		"user-id",
		"test@example.com",
		"https://example.com/reset",
	)
	assert.Equal(t, "user.password_reset_requested", event.EventName())
	assert.Equal(t, "user-id", event.AggregateID())
	assert.Equal(t, "v1", event.EventVersion())
	assert.Equal(t, "test@example.com", event.UserEmail)
	assert.Equal(t, "https://example.com/reset", event.ResetURL)
	assert.False(t, event.OccurredAt().IsZero(), "Expected OccurredAt to not be zero value")
	assert.True(t, event.OccurredAt().Before(time.Now().UTC()), "Expected OccurredAt to be recent")
}

func TestPasswordResetRequestedEventImplementsIEvent(t *testing.T) {
	event := NewPasswordResetRequestedEvent(
		"user-id",
		"test@example.com",
		"https://example.com/reset",
	)
	var iEvent IEvent = event
	assert.Equal(t, "user.password_reset_requested", iEvent.EventName())
	assert.Equal(t, "user-id", iEvent.AggregateID())
	assert.Equal(t, "v1", iEvent.EventVersion())
	assert.Equal(t, "test@example.com", event.UserEmail)
}
