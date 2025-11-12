// Package event
package event

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type EventFactory func() Event

type Registry struct {
	factories map[string]EventFactory
	mu        sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]EventFactory),
	}
}

func (r *Registry) Register(eventName string, factory EventFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[eventName]; exists {
		log.Printf("event already registered: %s", eventName)
	}

	r.factories[eventName] = factory
}

func (r *Registry) Unmarshal(data []byte) (Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var partial struct {
		Name string `json:"Name"`
	}
	if err := json.Unmarshal(data, &partial); err != nil {
		return nil, fmt.Errorf("failed to extract event name: %w", err)
	}
	if partial.Name == "" {
		return nil, fmt.Errorf("event name is empty")
	}
	factory, exists := r.factories[partial.Name]
	if !exists {
		return nil, fmt.Errorf("unknown event type: %s (not registered in registry)", partial.Name)
	}
	evt := factory()
	if err := json.Unmarshal(data, evt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event %s: %w", partial.Name, err)
	}
	return evt, nil
}

func (r *Registry) IsRegistered(eventName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.factories[eventName]
	return exists
}

func (r *Registry) RegisteredEvents() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	events := make([]string, 0, len(r.factories))
	for name := range r.factories {
		events = append(events, name)
	}
	return events
}

func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories = make(map[string]EventFactory)
}

func RegisterAll(registry *Registry) {
	registry.Register("user.password_reset_requested", func() Event {
		return &PasswordResetRequestedEvent{}
	})
}
