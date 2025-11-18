// Package event
package event

import "time"

type IEvent interface {
	EventName() string
	OccurredAt() time.Time
	EventVersion() string
	GetAggregateID() string
	Metadata() map[string]any
}

type Event struct {
	Name        string
	Timestamp   time.Time
	Version     string
	AggregateID string
	Meta        map[string]any
}

func (e Event) EventName() string {
	return e.Name
}

func (e Event) OccurredAt() time.Time {
	return e.Timestamp
}

func (e Event) EventVersion() string {
	return e.Version
}

func (e Event) GetAggregateID() string {
	return e.AggregateID
}

func (e Event) Metadata() map[string]any {
	if e.Meta == nil {
		return make(map[string]any)
	}
	return e.Meta
}

func NewEvent(name, aggregateID, version string) Event {
	return Event{
		Name:        name,
		Timestamp:   time.Now().UTC(),
		Version:     version,
		AggregateID: aggregateID,
		Meta:        make(map[string]any),
	}
}
