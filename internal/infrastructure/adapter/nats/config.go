// Package nats
package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

func ConnectAndSetup(cfg config.NATSConfig) (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := nats.Connect(
		cfg.URL,
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      cfg.StreamName,
		Subjects:  []string{"users.>"},
		Retention: nats.WorkQueuePolicy,
		MaxAge:    7 * 24 * time.Hour, // 7 dias
		Storage:   nats.FileStorage,
	})
	if err != nil && !isStreamExistsError(err) {
		nc.Close()
		return nil, nil, fmt.Errorf("failed to setup main stream: %w", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "EVENT_DLQ",
		Subjects:  []string{"user.>"},
		Retention: nats.LimitsPolicy,
		MaxAge:    30 * 24 * time.Hour, // 30 dias
		Storage:   nats.FileStorage,
	})
	if err != nil && !isStreamExistsError(err) {
		nc.Close()
		return nil, nil, fmt.Errorf("failed to setup DLQ stream: %w", err)
	}
	return nc, js, nil
}

func isStreamExistsError(err error) bool {
	if err, ok := err.(*nats.APIError); ok {
		return err.ErrorCode == 10058
	}
	return false
}
