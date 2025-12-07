// Package gateway_test
package gateway_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/test/spy"
)

// TestEvent implementa a interface IEvent para testes
type TestEvent struct {
	name        string
	timestamp   time.Time
	version     string
	aggregateID string
}

func (e TestEvent) EventName() string {
	return e.name
}

func (e TestEvent) OccurredAt() time.Time {
	return e.timestamp
}

func (e TestEvent) EventVersion() string {
	return e.version
}

func (e TestEvent) AggregateID() string {
	return e.aggregateID
}

func TestPublisherInterface(t *testing.T) {
	// Teste para garantir que o Publisher spy implementa a interface corretamente
	var publisher gateway.Publisher = &spy.Publisher{}
	
	// Verifica se é possível atribuir um spy.Publisher a um gateway.Publisher
	if publisher == nil {
		t.Error("spy.Publisher não implementa gateway.Publisher")
	}
}

func TestSpyPublisherPublish(t *testing.T) {
	ctx := context.Background()
	testEvent := TestEvent{
		name:        "test.event",
		timestamp:   time.Now(),
		version:     "1.0",
		aggregateID: "test-id",
	}

	spyPublisher := &spy.Publisher{}

	// Testa publicação com sucesso
	err := spyPublisher.Publish(ctx, testEvent)
	if err != spyPublisher.Err {
		t.Errorf("Esperava erro %v, mas obteve %v", spyPublisher.Err, err)
	}

	// Verifica se o evento foi armazenado corretamente
	if spyPublisher.Event == nil {
		t.Error("Evento não foi armazenado pelo spy publisher")
	} else if spyPublisher.Event.EventName() != testEvent.EventName() {
		t.Errorf("Nome do evento incorreto. Esperava %s, obteve %s", testEvent.EventName(), spyPublisher.Event.EventName())
	}
}

func TestSpyPublisherPublishBatch(t *testing.T) {
	ctx := context.Background()
	testEvents := []event.IEvent{
		TestEvent{
			name:        "test.event.1",
			timestamp:   time.Now(),
			version:     "1.0",
			aggregateID: "test-id-1",
		},
		TestEvent{
			name:        "test.event.2",
			timestamp:   time.Now(),
			version:     "1.0",
			aggregateID: "test-id-2",
		},
	}

	spyPublisher := &spy.Publisher{}
	
	// Testa publicação em lote com sucesso
	err := spyPublisher.PublishBatch(ctx, testEvents)
	if err != spyPublisher.Err {
		t.Errorf("Esperava erro %v, mas obteve %v", spyPublisher.Err, err)
	}
}

func TestSpyPublisherClose(t *testing.T) {
	spyPublisher := &spy.Publisher{}
	
	// Testa fechamento com sucesso
	err := spyPublisher.Close()
	if err != spyPublisher.Err {
		t.Errorf("Esperava erro %v, mas obteve %v", spyPublisher.Err, err)
	}
}

func TestSpyPublisherWithError(t *testing.T) {
	ctx := context.Background()
	testEvent := TestEvent{
		name:        "test.event",
		timestamp:   time.Now(),
		version:     "1.0",
		aggregateID: "test-id",
	}

	expectedError := errors.New("publisher error")
	spyPublisher := &spy.Publisher{Err: expectedError}
	
	// Testa publicação com erro
	err := spyPublisher.Publish(ctx, testEvent)
	if err != expectedError {
		t.Errorf("Esperava erro %v, mas obteve %v", expectedError, err)
	}
}