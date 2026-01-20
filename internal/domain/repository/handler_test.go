package repository

import (
	"context"
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/stretchr/testify/assert"
)

type EventSpy struct{}

func (e EventSpy) EventName() string {
	return "mock.event"
}

func (e EventSpy) OccurredAt() time.Time {
	return time.Now()
}

func (e EventSpy) EventVersion() string {
	return "1.0"
}

func (e EventSpy) AggregateID() string {
	return "mock-aggregate-id"
}

type MockHandler struct{}

func (e *MockHandler) Handle(ctx context.Context, event event.IEvent) error {
	return nil
}

func (e *MockHandler) CanHandle(event event.IEvent) bool {
	return true
}

func TestHandlerInterface(t *testing.T) {
	ctx := context.Background()
	sut := &MockHandler{}
	mockEvent := EventSpy{}
	err := sut.Handle(ctx, mockEvent)
	assert.NoError(t, err)
	canHandle := sut.CanHandle(mockEvent)
	assert.True(t, canHandle)
}

