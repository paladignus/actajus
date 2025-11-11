// Package event
package event

import (
	"fmt"
	"sync"
)

type EventFactory func() Event

type EventRegistry struct {
	mu        sync.RWMutex
	factories map[string]EventFactory
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		factories: make(map[string]EventFactory),
	}
}

func (r *EventRegistry) Register(eventName string, factory EventFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.factories[eventName]; exists {
		return fmt.Errorf("event type already registered: %s", eventName)
	}
	r.factories[eventName] = factory
	return nil
}

func (r *EventRegistry) CreateEventByName(eventName string) (Event, error) {
	r.mu.RLock()
	factory, exists := r.factories[eventName]
	r.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("unknown event type: %s", eventName)
	}
	return factory(), nil
}

func (r *EventRegistry) GetRegisteredEventNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}
