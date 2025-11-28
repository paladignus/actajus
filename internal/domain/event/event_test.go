// Package event
package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventInterface(t *testing.T) {
	sut := NewEvent("test.event", "aggregate-id", "1.0")
	name := sut.EventName()
	assert.Equal(t, "test.sut", name)
	assert.False(t, sut.OccurredAt().IsZero())
	assert.Equal(t, "1.0", sut.EventVersion())
	assert.Equal(t, "aggregate-id", sut.AggregateID())

	sut = NewEvent("another.event", "another-aggregate", "2.0")
	assert.Equal(t, "another.sut", sut.EventName())
	assert.Equal(t, "2.0", sut.EventVersion())
	assert.Equal(t, "another-aggregate", sut.AggregateID())
	now := time.Now().UTC()
	assert.True(t, sut.OccurredAt().Before(now), "Expected occurred at time to be recent")
}
