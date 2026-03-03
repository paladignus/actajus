// Package dispatcher
package dispatcher

import "time"

// EventDispatcher define um contrato para publicação de eventos de domínio
type EventDispatcher interface {
	// Publish publica um evento para todos os listeners registrados
	Publish(event any)
	// Register registra um listener para um tipo específico de evento
	Register(eventType string, listener EventListener)
}

// EventListener define um contrato para listeners de eventos
type EventListener interface {
	// Handle processa um evento
	Handle(event any)
	// EventType retorna o tipo de evento que este listener processa
	EventType() string
}

// EventHandler é uma função que processa um evento
type EventHandler func(event any)

// SimpleEventDispatcher é uma implementação simples de EventDispatcher
type SimpleEventDispatcher struct {
	listeners map[string][]EventHandler
}

// NewSimpleEventDispatcher cria um novo EventDispatcher
func NewSimpleEventDispatcher() *SimpleEventDispatcher {
	return &SimpleEventDispatcher{
		listeners: make(map[string][]EventHandler),
	}
}

// Publish publica um evento para todos os listeners registrados
func (d *SimpleEventDispatcher) Publish(event any) {
	eventType := getEventType(event)
	handlers, ok := d.listeners[eventType]
	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(event)
	}
}

// Register registra um handler para um tipo específico de evento
func (d *SimpleEventDispatcher) Register(eventType string, handler EventHandler) {
	if d.listeners == nil {
		d.listeners = make(map[string][]EventHandler)
	}
	d.listeners[eventType] = append(d.listeners[eventType], handler)
}

// getEventType retorna o tipo de evento como string
func getEventType(event any) string {
	switch event.(type) {
	case RoleAssignedEvent:
		return "role_assigned"
	case RoleRemovedEvent:
		return "role_removed"
	case PermissionGrantedEvent:
		return "permission_granted"
	case PermissionRevokedEvent:
		return "permission_revoked"
	default:
		return "unknown"
	}
}

// Eventos de domínio do módulo Identity
type RoleAssignedEvent struct {
	ID         int64
	IDRole     int16
	AssignedBy int64
	AssignedAt time.Time
}

type RoleRemovedEvent struct {
	IDUser     int64
	IDRole     int16
	RemovedAt  time.Time
}

type PermissionGrantedEvent struct {
	IDRole       int16
	IDPermission int16
	GrantedAt    time.Time
}

type PermissionRevokedEvent struct {
	IDRole       int16
	IDPermission int16
	RevokedAt    time.Time
}
