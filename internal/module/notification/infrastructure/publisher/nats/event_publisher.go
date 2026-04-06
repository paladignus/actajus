// Package nats
package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/paladignus/actajus/internal/module/notification/application/service"
	"github.com/paladignus/actajus/internal/module/notification/domain"
)

// JetStreamEventPublisher é uma implementação de EventPublisher usando NATS JetStream
type JetStreamEventPublisher struct {
	js        jetstream.JetStream
	stream    string
	subjectFn func(eventType string) string
}

// NewJetStreamEventPublisher cria um novo publicador de eventos usando NATS JetStream
func NewJetStreamEventPublisher(
	js jetstream.JetStream,
	stream string,
	subjectFn func(eventType string) string,
) *JetStreamEventPublisher {
	return &JetStreamEventPublisher{
		js:        js,
		stream:    stream,
		subjectFn: subjectFn,
	}
}

// Publish publica um evento de domínio
func (p *JetStreamEventPublisher) Publish(ctx context.Context, event domain.IEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	subject := p.subjectFn(event.EventType())
	_, err = p.js.Publish(ctx, subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish event %s: %w", event.EventType(), err)
	}

	return nil
}

// PublishBatch publica múltiplos eventos de domínio
func (p *JetStreamEventPublisher) PublishBatch(ctx context.Context, events []domain.IEvent) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// DefaultSubjectFn é a função padrão para mapear tipo de evento para subject NATS
func DefaultSubjectFn(eventType string) string {
	// Ex: "user.created" → "notification.event.user.created"
	return fmt.Sprintf("notification.event.%s", eventType)
}

// EnsureStream garante que o stream de eventos existe
func EnsureStream(ctx context.Context, js jetstream.JetStream, streamName string) error {
	_, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  []string{"notification.event.>"},
		Replicas:  1,
		Retention: jetstream.WorkQueuePolicy,
	})
	if err != nil {
		// Stream já pode existir
		return nil
	}
	return nil
}

// Compile-time check
var _ service.EventPublisher = (*JetStreamEventPublisher)(nil)
