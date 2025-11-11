// Package event
package event

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// EventFactory é uma função que cria uma nova instância de um evento.
// Retorna um ponteiro para o tipo concreto do evento.
type EventFactory func() Event

// Registry mantém o registro de todos os tipos de eventos do sistema.
// Permite deserialização automática de JSON para o tipo correto.
//
// Funcionamento:
//  1. Cada tipo de evento é registrado com seu nome
//  2. Ao deserializar, o Registry identifica o tipo pelo nome
//  3. Cria uma instância do tipo correto
//  4. Faz unmarshal do JSON nessa instância
type Registry struct {
	factories map[string]EventFactory
	mu        sync.RWMutex
}

// NewRegistry cria um novo registro de eventos vazio.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]EventFactory),
	}
}

// Register registra um novo tipo de evento no registry.
//
// Parâmetros:
//   - eventName: nome único do evento (ex: "user.password_reset_requested")
//   - factory: função que cria uma nova instância deste tipo de evento
//
// Exemplo de uso:
//
//	registry.Register("user.password_reset_requested", func() Event {
//	    return &PasswordResetRequestedEvent{}
//	})
func (r *Registry) Register(eventName string, factory EventFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[eventName]; exists {
		log.Printf("event already registered: %s", eventName)
		// Permite re-registro (útil em testes), mas loga aviso
		// Em produção, você pode querer retornar erro aqui
		// panic(fmt.Sprintf("event already registered: %s", eventName))
	}

	r.factories[eventName] = factory
}

// Unmarshal deserializa dados JSON para o tipo de evento correto.
//
// Parâmetros:
//   - data: bytes JSON do evento serializado
//
// Retorna:
//   - Event: instância do tipo correto do evento
//   - error: se não conseguir deserializar ou tipo não estiver registrado
//
// O que faz:
//  1. Faz unmarshal parcial para extrair o EventName
//  2. Busca a factory correspondente no registry
//  3. Cria nova instância do tipo correto
//  4. Faz unmarshal completo nessa instância
//  5. Retorna o evento tipado
func (r *Registry) Unmarshal(data []byte) (Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Primeira passagem: extrai apenas o nome do evento
	var partial struct {
		Name string `json:"Name"`
	}

	if err := json.Unmarshal(data, &partial); err != nil {
		return nil, fmt.Errorf("failed to extract event name: %w", err)
	}

	if partial.Name == "" {
		return nil, fmt.Errorf("event name is empty")
	}

	// Busca a factory para este tipo de evento
	factory, exists := r.factories[partial.Name]
	if !exists {
		return nil, fmt.Errorf("unknown event type: %s (not registered in registry)", partial.Name)
	}

	// Cria nova instância do tipo correto
	evt := factory()

	// Segunda passagem: unmarshal completo no tipo correto
	if err := json.Unmarshal(data, evt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event %s: %w", partial.Name, err)
	}

	return evt, nil
}

// IsRegistered verifica se um tipo de evento está registrado.
func (r *Registry) IsRegistered(eventName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.factories[eventName]
	return exists
}

// RegisteredEvents retorna a lista de todos os eventos registrados.
// Útil para debugging e validação.
func (r *Registry) RegisteredEvents() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]string, 0, len(r.factories))
	for name := range r.factories {
		events = append(events, name)
	}
	return events
}

// Clear remove todos os eventos registrados.
// Útil principalmente em testes.
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories = make(map[string]EventFactory)
}

// GlobalRegistry é uma instância global do registry para conveniência.
// Você pode usar esta instância ou criar suas próprias instâncias.
var GlobalRegistry = NewRegistry()

// Register é um helper que registra no GlobalRegistry.
// Equivalente a GlobalRegistry.Register(eventName, factory)
func Register(eventName string, factory EventFactory) {
	GlobalRegistry.Register(eventName, factory)
}

// Unmarshal é um helper que deserializa usando GlobalRegistry.
// Equivalente a GlobalRegistry.Unmarshal(data)
func Unmarshal(data []byte) (Event, error) {
	return GlobalRegistry.Unmarshal(data)
}

// type EventFactory func() Event
//
// type EventRegistry struct {
// 	mu        sync.RWMutex
// 	factories map[string]EventFactory
// }
//
// func NewEventRegistry() *EventRegistry {
// 	return &EventRegistry{
// 		factories: make(map[string]EventFactory),
// 	}
// }
//
// func (r *EventRegistry) Register(eventName string, factory EventFactory) error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	if _, exists := r.factories[eventName]; exists {
// 		return fmt.Errorf("event type already registered: %s", eventName)
// 	}
// 	r.factories[eventName] = factory
// 	return nil
// }
//
// func (r *EventRegistry) CreateEventByName(eventName string) (Event, error) {
// 	r.mu.RLock()
// 	factory, exists := r.factories[eventName]
// 	r.mu.RUnlock()
// 	if !exists {
// 		return nil, fmt.Errorf("unknown event type: %s", eventName)
// 	}
// 	return factory(), nil
// }
//
// func (r *EventRegistry) GetRegisteredEventNames() []string {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()
// 	names := make([]string, 0, len(r.factories))
// 	for name := range r.factories {
// 		names = append(names, name)
// 	}
// 	return names
// }
