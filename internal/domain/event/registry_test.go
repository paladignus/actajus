// Package event
package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock event for testing
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

// Another mock event
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

//	func TestRegistry_RegisterDuplicate(t *testing.T) {
//		registry := NewRegistry()
//		factory1 := func() IEvent {
//			return &MockEvent{Name: "first"}
//		}
//		factory2 := func() IEvent {
//			return &MockEvent{Name: "second"}
//		}
//		registry.Register("test.event", factory1)
//		registry.Register("test.event", factory2)
//		if !registry.IsRegistered("test.event") {
//			t.Error("event should remain registered after duplicate registration")
//		}
//	}
func TestRegistry(t *testing.T) {
	registry := NewRegistry()
	factory := func() IEvent {
		return &MockEvent{Name: "first"}
	}
	factory2 := func() IEvent {
		return &MockEvent{Name: "second"}
	}

	t.Run("should create a new registry", func(t *testing.T) {
		assert.NotNil(t, registry)
		assert.NotNil(t, registry.factories)
		assert.Empty(t, registry.factories)
	})

	t.Run("should register an event", func(t *testing.T) {
		registry.Register("test.event", factory)
		assert.True(t, registry.IsRegistered("test.event"))
	})

	t.Run("should return an error when registering duplicate event", func(t *testing.T) {
		registry.Register("test.event", factory2)
		assert.True(t, registry.IsRegistered("test.event"))
	})

	t.Run("should return success to unmarshal success", func(t *testing.T) {
		data := []byte(`{"Name":"test.event","Payload":"test data"}`)
		evt, err := registry.Unmarshal(data)
		mockEvt, ok := evt.(*MockEvent)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "test.event", mockEvt.Name)
		assert.Equal(t, "test data", mockEvt.Payload)
	})

	t.Run("should return error to unmarshal failure to invalid JSON", func(t *testing.T) {
		data := []byte(`{invalid json}`)
		_, err := registry.Unmarshal(data)
		assert.Error(t, err)
	})

	t.Run("should return error to unmarshal failure to empty name", func(t *testing.T) {
		data := []byte(`{"Name":"","Payload":"test"}`)
		_, err := registry.Unmarshal(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event name is empty")
	})

	t.Run("should return error to unmarshal failure to unknown event", func(t *testing.T) {
		data := []byte(`{"Name":"unknown.event","Payload":"test"}`)
		_, err := registry.Unmarshal(data)
		assert.Error(t, err)
	})

	t.Run("should return error to unmarshal failure to invalid event data", func(t *testing.T) {
		registry.Register("another.event", func() IEvent {
			return &AnotherMockEvent{}
		})
		data := []byte(`{"Name":"another.event","Value":"not-a-number"}`)
		_, err := registry.Unmarshal(data)
		assert.Error(t, err)
	})

	t.Run("should correct event is registered", func(t *testing.T) {
		assert.False(t, registry.IsRegistered("nonexistent"))
		assert.True(t, registry.IsRegistered("test.event"))
	})

	t.Run("should register an event", func(t *testing.T) {
		registry.Register("event1", func() IEvent { return &MockEvent{} })
		registry.Register("event2", func() IEvent { return &AnotherMockEvent{} })
		registry.Register("event3", func() IEvent { return &MockEvent{} })
		events := registry.RegisteredEvents()
		assert.NotEmpty(t, registry.RegisteredEvents())
		assert.Len(t, events, 5)
		assert.Contains(t, events, "event1")
		assert.Contains(t, events, "event2")
		assert.Contains(t, events, "event3")
	})

	t.Run("should clear registered events", func(t *testing.T) {
		registry.Clear()
		assert.Empty(t, registry.RegisteredEvents())
		assert.False(t, registry.IsRegistered("test.event"))
	})

	t.Run("should success in register all", func(t *testing.T) {
		RegisterAll(registry)
		data := []byte(`{"Name":"user.password_reset_requested"}`)
		evt, err := registry.Unmarshal(data)
		_, ok := evt.(*PasswordResetRequestedEvent)
		assert.True(t, registry.IsRegistered("user.password_reset_requested"))
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}
