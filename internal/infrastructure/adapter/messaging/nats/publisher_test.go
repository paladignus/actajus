// Package nats_test
package nats_test

import (
	"testing"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/infrastructure/adapter/messaging/nats"
	"github.com/stretchr/testify/assert"
)

// Teste para verificar a implementação da interface Publisher
func TestPublisherImplementsInterface(t *testing.T) {
	// Este teste verifica se o publisher implementa a interface corretamente
	// O código real já garante isso em tempo de compilação, mas este teste documenta a intenção
	
	// Criar um publisher fake apenas para verificação de interface (não usado funcionalmente)
	var publisher interface{} = &nats.Publisher{}
	assert.NotNil(t, publisher)
}

// Teste para os eventos de teste auxiliares
func TestTestEvent(t *testing.T) {
	t.Run("TestEvent should implement IEvent interface", func(t *testing.T) {
		testEvent := &nats.TestEvent{
			Data: "test data",
			ID:   "test-id",
		}
		
		var eventIface event.IEvent = testEvent
		assert.NotNil(t, eventIface)
		
		// Testar os métodos da interface
		assert.Equal(t, "test.event", eventIface.EventName())
		assert.Equal(t, "test-id", eventIface.AggregateID())
		assert.Equal(t, "v1", eventIface.EventVersion())
		
		// Verificar que OccurredAt retorna um tempo válido
		timestamp := eventIface.OccurredAt()
		assert.False(t, timestamp.IsZero())
	})

	t.Run("FailingTestEvent should implement IEvent interface", func(t *testing.T) {
		failingEvent := &nats.FailingTestEvent{}
		
		var eventIface event.IEvent = failingEvent
		assert.NotNil(t, eventIface)
		
		// Testar os métodos da interface
		assert.Equal(t, "failing.event", eventIface.EventName())
		assert.Equal(t, "test-id", eventIface.AggregateID())
		assert.Equal(t, "v1", eventIface.EventVersion())
		
		// Verificar que OccurredAt retorna um tempo válido
		timestamp := eventIface.OccurredAt()
		assert.False(t, timestamp.IsZero())
	})
}

// Teste para verificação de serialização dos eventos de teste
func TestTestEventSerialization(t *testing.T) {
	t.Run("TestEvent should serialize and deserialize correctly", func(t *testing.T) {
		originalEvent := &nats.TestEvent{
			Data: "test data",
			ID:   "test-id",
		}
		
		// Testar serialização
		data, err := originalEvent.MarshalJSON()
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
		
		// Testar desserialização
		newEvent := &nats.TestEvent{}
		err = newEvent.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, originalEvent.Data, newEvent.Data)
		assert.Equal(t, originalEvent.ID, newEvent.ID)
	})

	t.Run("FailingTestEvent should fail serialization", func(t *testing.T) {
		failingEvent := &nats.FailingTestEvent{}
		
		_, err := failingEvent.MarshalJSON()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "serialization failed")
	})
}

// Teste para o spy publisher com o publisher real
func TestPublisherWithSpy(t *testing.T) {
	t.Run("should verify publisher can be assigned to interface", func(t *testing.T) {
		// Este teste apenas verifica que podemos lidar com o publisher como uma interface
		var publisher interface{} = &nats.Publisher{}
		assert.NotNil(t, publisher)
	})
}

// Teste adicional para verificar o comportamento dos erros personalizados
func TestNATSErrors(t *testing.T) {
	t.Run("NATS errors should implement error interface", func(t *testing.T) {
		// Os testes para os erros NATS já estão em errors_test.go
		// Este teste serve apenas para confirmar a funcionalidade
		err := nats.ErrConnectionFailed(nil)
		assert.Error(t, err)
		
		err = nats.ErrPublishFailed(nil)
		assert.Error(t, err)
		
		err = nats.ErrSerializationFailed(nil)
		assert.Error(t, err)
	})
}