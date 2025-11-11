// Package event
package event

import "time"

type Event interface {
	EventName() string
	OccurredAt() time.Time
	EventVersion() string
	GetAggregateID() string
	Metadata() map[string]any
}

type BaseEvent struct {
	Name        string
	Timestamp   time.Time
	Version     string
	AggregateID string
	Meta        map[string]any
}

func (e BaseEvent) EventName() string {
	return e.Name
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.Timestamp
}

func (e BaseEvent) EventVersion() string {
	return e.Version
}

func (e BaseEvent) GetAggregateID() string {
	return e.AggregateID
}

func (e BaseEvent) Metadata() map[string]any {
	if e.Meta == nil {
		return make(map[string]any)
	}
	return e.Meta
}

func NewBaseEvent(name, aggregateID, version string) BaseEvent {
	return BaseEvent{
		Name:        name,
		Timestamp:   time.Now().UTC(),
		Version:     version,
		AggregateID: aggregateID,
		Meta:        make(map[string]any),
	}
}
