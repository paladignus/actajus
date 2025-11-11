// Package nats
package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/paladignus/actajus/internal/domain/event"
)

type Publisher struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	config *Config
}

func NewPublisher(config *Config) (*Publisher, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	conn, err := nats.Connect(
		config.URL,
		nats.Timeout(config.ConnectionTimeout),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1), // Reconecta infinitamente
		nats.ReconnectWait(2*time.Second),
		// Callbacks para observabilidade
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				fmt.Printf("[NATS] Disconnected: %v\n", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("[NATS] Reconnected to %s\n", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		return nil, ErrConnectionFailed(err)
	}
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, ErrConnectionFailed(err)
	}
	publisher := &Publisher{
		conn:   conn,
		js:     js,
		config: config,
	}
	if err := publisher.ensureStream(); err != nil {
		conn.Close()
		return nil, err
	}
	return publisher, nil
}

func (p *Publisher) ensureStream() error {
	streamConfig := &nats.StreamConfig{
		Name:     p.config.StreamName,
		Subjects: p.config.Subjects,
		MaxAge:   p.config.MaxAge,
		MaxBytes: p.config.MaxBytes,
		Replicas: p.config.Replicas,
		Storage:  nats.FileStorage,
	}
	_, err := p.js.AddStream(streamConfig)
	if err != nil {
		// Se o stream já existe, tenta atualizar
		_, err = p.js.UpdateStream(streamConfig)
		if err != nil {
			return ErrStreamCreationFailed(err)
		}
	}
	return nil
}

func (p *Publisher) Publish(ctx context.Context, evt event.Event) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return ErrSerializationFailed(err)
	}
	subject := fmt.Sprintf("events.%s", evt.EventName())
	log.Printf("[NATS] Publishing event: %s, %v", subject, evt)
	_, err = p.js.Publish(subject, data, nats.Context(ctx))
	if err != nil {
		return ErrPublishFailed(err)
	}
	return nil
}

func (p *Publisher) PublishBatch(ctx context.Context, events []event.Event) error {
	var errors []error
	for _, evt := range events {
		if err := p.Publish(ctx, evt); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("failed to publish %d events: %v", len(errors), errors)
	}
	return nil
}

func (p *Publisher) Close() error {
	if p.conn != nil {
		// Aguarda mensagens pendentes serem enviadas
		p.conn.Flush()
		// Fecha a conexão
		p.conn.Close()
	}
	return nil
}
