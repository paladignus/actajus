// Package domain
package domain

import "time"

// IEvent define o contrato para eventos de domínio
type IEvent interface {
	// EventType retorna o tipo do evento (ex: "user.created", "company.deleted")
	EventType() string
	// OccurredAt retorna quando o evento ocorreu
	OccurredAt() time.Time
	// AggregateID retorna o ID da agregação raiz que originou o evento
	AggregateID() string
	// Version retorna a versão do schema do evento
	Version() string
}

// Event é a implementação base para eventos de domínio
type Event struct {
	eventType   string
	occurredAt  time.Time
	aggregateID string
	version     string
	metadata    map[string]any
}

// EventType retorna o tipo do evento
func (e Event) EventType() string {
	return e.eventType
}

// OccurredAt retorna quando o evento ocorreu
func (e Event) OccurredAt() time.Time {
	return e.occurredAt
}

// AggregateID retorna o ID da agregação raiz
func (e Event) AggregateID() string {
	return e.aggregateID
}

// Version retorna a versão do schema
func (e Event) Version() string {
	return e.version
}

// Metadata retorna metadados do evento
func (e Event) Metadata() map[string]any {
	return e.metadata
}

// NewEvent cria um novo evento de domínio
func NewEvent(eventType, aggregateID, version string, metadata map[string]any) Event {
	return Event{
		eventType:   eventType,
		occurredAt:  time.Now().UTC(),
		aggregateID: aggregateID,
		version:     version,
		metadata:    metadata,
	}
}
