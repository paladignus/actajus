// Package event
package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasswordResetRequestedEvent(t *testing.T) {
	sut := NewPasswordResetRequestedEvent(
		"user-id",
		"test@example.com",
		"https://example.com/reset",
	)
	assert.Equal(t, "user.password_reset_requested", sut.EventName())
	assert.Equal(t, "user-id", sut.AggregateID())
	assert.Equal(t, "v1", sut.EventVersion())
	assert.Equal(t, "test@example.com", sut.UserEmail)
	assert.Equal(t, "https://example.com/reset", sut.ResetURL)
	assert.False(t, sut.OccurredAt().IsZero(), "Expected OccurredAt to not be zero value")
	assert.True(t, sut.OccurredAt().Before(time.Now().UTC()), "Expected OccurredAt to be recent")
}

func TestPasswordResetRequestedEventImplementsIEvent(t *testing.T) {
	sut := NewPasswordResetRequestedEvent(
		"user-id",
		"test@example.com",
		"https://example.com/reset",
	)
	var iEvent IEvent = sut
	assert.Equal(t, "user.password_reset_requested", iEvent.EventName())
	assert.Equal(t, "user-id", iEvent.AggregateID())
	assert.Equal(t, "v1", iEvent.EventVersion())
	assert.Equal(t, "test@example.com", sut.UserEmail)
}
