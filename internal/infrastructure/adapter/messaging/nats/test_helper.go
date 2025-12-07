package nats

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/paladignus/actajus/internal/domain/event"
)

// TestEvent is a test event for testing purposes
type TestEvent struct {
	Data string `json:"data"`
	ID   string `json:"id"`
}

func (m *TestEvent) EventName() string {
	return "test.event"
}

func (m *TestEvent) OccurredAt() time.Time {
	return time.Now()
}

func (m *TestEvent) EventVersion() string {
	return "v1"
}

func (m *TestEvent) AggregateID() string {
	return "test-id"
}

func (m *TestEvent) MarshalJSON() ([]byte, error) {
	// Use uma struct anônima ou uma cópia para evitar o loop infinito
	return json.Marshal(struct {
		Data string `json:"data"`
		ID   string `json:"id"`
	}{
		Data: m.Data,
		ID:   m.ID,
	})
}

func (m *TestEvent) UnmarshalJSON(data []byte) error {
	aux := struct {
		Data string `json:"data"`
		ID   string `json:"id"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	m.Data = aux.Data
	m.ID = aux.ID
	return nil
}

// FailingTestEvent is a test event that fails to serialize
type FailingTestEvent struct{}

func (f *FailingTestEvent) EventName() string {
	return "failing.event"
}

func (f *FailingTestEvent) OccurredAt() time.Time {
	return time.Now()
}

func (f *FailingTestEvent) EventVersion() string {
	return "v1"
}

func (f *FailingTestEvent) AggregateID() string {
	return "test-id"
}

func (f *FailingTestEvent) MarshalJSON() ([]byte, error) {
	return nil, errors.New("serialization failed")
}

func (f *FailingTestEvent) UnmarshalJSON(data []byte) error {
	return nil
}

// ensure interfaces are implemented
var (
	_ event.IEvent = (*TestEvent)(nil)
	_ event.IEvent = (*FailingTestEvent)(nil)
)