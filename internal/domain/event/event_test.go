// Package event
package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventInterface(t *testing.T) {
	event := NewEvent("test.event", "aggregate-id", "1.0")
	name := event.EventName()
	assert.Equal(t, "test.event", name)
	assert.False(t, event.OccurredAt().IsZero())
	assert.Equal(t, "1.0", event.EventVersion())
	assert.Equal(t, "aggregate-id", event.AggregateID())

	event = NewEvent("another.event", "another-aggregate", "2.0")
	assert.Equal(t, "another.event", event.EventName())
	assert.Equal(t, "2.0", event.EventVersion())
	assert.Equal(t, "another-aggregate", event.AggregateID())
	now := time.Now().UTC()
	assert.True(t, event.OccurredAt().Before(now), "Expected occurred at time to be recent")
}
