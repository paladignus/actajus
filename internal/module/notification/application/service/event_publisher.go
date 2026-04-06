// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/module/notification/domain"
)

// EventPublisher define o contrato para publicação de eventos de domínio
type EventPublisher interface {
	// Publish publica um evento de domínio
	Publish(ctx context.Context, event domain.IEvent) error
	// PublishBatch publica múltiplos eventos de domínio
	PublishBatch(ctx context.Context, events []domain.IEvent) error
}

// EventSubscriber define o contrato para subscrição de eventos de domínio
type EventSubscriber interface {
	// Subscribe se inscreve para receber eventos de um tipo específico
	Subscribe(eventType string, handler EventHandler) error
	// Unsubscribe cancela uma inscrição
	Unsubscribe(eventType string) error
}

// EventHandler é uma função que processa eventos de domínio
type EventHandler func(ctx context.Context, event domain.IEvent) error
