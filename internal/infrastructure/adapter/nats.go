// Package adapter
package adapter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

const (
	streamName       = "EVENTS"
	dlqStreamName    = "EVENTS_DLQ"
	durablePrefix    = "welcome-email"
	JSStreamExistErr = 10058 // NATS error code for "stream name already in use"
)

type natsAdapter struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewNATSAdapter(config config.NATSConfig) (*natsAdapter, error) {
	nc, err := nats.Connect(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	// Setup stream principal
	err = ensureStream(js, streamName, []string{"user.>"})
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to setup main stream: %w", err)
	}

	// Setup DLQ stream
	err = ensureStream(js, dlqStreamName, []string{"dlq.user.>"})
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to setup DLQ stream: %w", err)
	}

	return &natsAdapter{nc: nc, js: js}, nil
}

// Publish implementa gateway.Publisher
func (n *natsAdapter) Publish(ctx context.Context, subject string, msg []byte) error {
	_, err := n.js.Publish(subject, msg)
	return err
}

// ensureStream cria um stream se ele não existir
func ensureStream(js nats.JetStreamContext, name string, subjects []string) error {
	_, err := js.AddStream(&nats.StreamConfig{
		Name:      name,
		Subjects:  subjects,
		Retention: nats.WorkQueuePolicy,
		MaxAge:    7 * 24 * time.Hour, // 7 dias
	})
	if err != nil {
		if isStreamExistsError(err) {
			return nil // já existe → OK
		}
		return err
	}
	return nil
}

// isStreamExistsError verifica se o erro é "stream already exists"
func isStreamExistsError(err error) bool {
	if apiErr, ok := err.(*nats.APIError); ok {
		return apiErr.ErrorCode == JSStreamExistErr
	}
	return false
}

func (n *natsAdapter) SubscribeWithRetry(
	ctx context.Context,
	subject string,
	durableName string,
	handler func(msg []byte) error,
) error {
	// Consumer com max 3 tentativas
	_, err := n.js.Subscribe(
		subject,
		func(m *nats.Msg) {
			info, err := m.Metadata()
			if err != nil {
				log.Printf("Failed to get message metadata: %v", err)
				m.Nak()
				return
			}

			if err := handler(m.Data); err != nil {
				log.Printf("Handler failed (attempt %d): %v", info.NumDelivered, err)

				// Se atingiu o limite de tentativas → enviar para DLQ
				if info.NumDelivered >= 3 {
					log.Printf("Moving message to DLQ after %d failures", info.NumDelivered)

					// Republish na DLQ com novo subject: "dlq.user.created"
					dlqSubject := "dlq." + subject
					if pubErr := n.Publish(context.Background(), dlqSubject, m.Data); pubErr != nil {
						log.Printf("Failed to publish to DLQ: %v", pubErr)
					}
					m.Ack() // ack na fila original para não reprocessar
				} else {
					// Retry com delay
					m.NakWithDelay(5 * time.Second)
				}
			} else {
				m.Ack()
			}
		},
		nats.Durable(durableName),
		nats.ManualAck(),
		nats.MaxDeliver(10), // alto limite, pois controlamos manualmente
		nats.AckWait(30*time.Second),
	)

	return err
}

func (n *natsAdapter) Close() {
	if n.nc != nil {
		n.nc.Close()
	}
}
