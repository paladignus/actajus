// Package event
package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockEvent struct {
	Name    string `json:"Name"`
	Payload string `json:"Payload"`
}

func (m *MockEvent) GetName() string {
	return m.Name
}

func (m *MockEvent) GetPayload() string {
	return m.Payload
}

func (m *MockEvent) EventName() string {
	return "mock.event"
}

func (m *MockEvent) AggregateID() string {
	return "mock-aggregate-id"
}

func (m *MockEvent) EventVersion() string {
	return "1.0"
}

func (m *MockEvent) OccurredAt() time.Time {
	return time.Now()
}

type AnotherMockEvent struct {
	Name  string `json:"Name"`
	Value int    `json:"Value"`
}

func (a *AnotherMockEvent) GetName() string {
	return a.Name
}

func (a *AnotherMockEvent) EventName() string {
	return "another.event"
}

func (a *AnotherMockEvent) AggregateID() string {
	return "another-aggregate"
}

func (a *AnotherMockEvent) EventVersion() string {
	return "1.0"
}

func (a *AnotherMockEvent) OccurredAt() time.Time {
	return time.Now()
}

func TestRegistry(t *testing.T) {
	sut := NewRegistry()
	factory := func() IEvent {
		return &MockEvent{Name: "first"}
	}
	factory2 := func() IEvent {
		return &MockEvent{Name: "second"}
	}

	t.Run("should create a new sut", func(t *testing.T) {
		assert.NotNil(t, sut)
		assert.NotNil(t, sut.factories)
		assert.Empty(t, sut.factories)
	})

	t.Run("should register an event", func(t *testing.T) {
		sut.Register("test.event", factory)
		assert.True(t, sut.IsRegistered("test.event"))
	})

	t.Run("should return an error when registering duplicate event", func(t *testing.T) {
		sut.Register("test.event", factory2)
		assert.True(t, sut.IsRegistered("test.event"))
	})

	t.Run("should return success to unmarshal success", func(t *testing.T) {
		data := []byte(`{"Name":"test.event","Payload":"test data"}`)
		evt, err := sut.Unmarshal(data)
		mockEvt, ok := evt.(*MockEvent)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "test.event", mockEvt.Name)
		assert.Equal(t, "test data", mockEvt.Payload)
	})

	t.Run("should return error to unmarshal failure to invalid JSON", func(t *testing.T) {
		data := []byte(`{invalid json}`)
		_, err := sut.Unmarshal(data)
		assert.Error(t, err)
	})

	t.Run("should return error to unmarshal failure to empty name", func(t *testing.T) {
		data := []byte(`{"Name":"","Payload":"test"}`)
		_, err := sut.Unmarshal(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event name is empty")
	})

	t.Run("should return error to unmarshal failure to unknown event", func(t *testing.T) {
		data := []byte(`{"Name":"unknown.event","Payload":"test"}`)
		_, err := sut.Unmarshal(data)
		assert.Error(t, err)
	})

	t.Run("should return error to unmarshal failure to invalid event data", func(t *testing.T) {
		sut.Register("another.event", func() IEvent {
			return &AnotherMockEvent{}
		})
		data := []byte(`{"Name":"another.event","Value":"not-a-number"}`)
		_, err := sut.Unmarshal(data)
		assert.Error(t, err)
	})

	t.Run("should correct event is registered", func(t *testing.T) {
		assert.False(t, sut.IsRegistered("nonexistent"))
		assert.True(t, sut.IsRegistered("test.event"))
	})

	t.Run("should register an event", func(t *testing.T) {
		sut.Register("event1", func() IEvent { return &MockEvent{} })
		sut.Register("event2", func() IEvent { return &AnotherMockEvent{} })
		sut.Register("event3", func() IEvent { return &MockEvent{} })
		events := sut.RegisteredEvents()
		assert.NotEmpty(t, sut.RegisteredEvents())
		assert.Len(t, events, 5)
		assert.Contains(t, events, "event1")
		assert.Contains(t, events, "event2")
		assert.Contains(t, events, "event3")
	})

	t.Run("should clear registered events", func(t *testing.T) {
		sut.Clear()
		assert.Empty(t, sut.RegisteredEvents())
		assert.False(t, sut.IsRegistered("test.event"))
	})

	t.Run("should success in register all", func(t *testing.T) {
		RegisterAll(sut)
		data := []byte(`{"Name":"user.password_reset_requested"}`)
		evt, err := sut.Unmarshal(data)
		_, ok := evt.(*PasswordResetRequestedEvent)
		assert.True(t, sut.IsRegistered("user.password_reset_requested"))
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}
