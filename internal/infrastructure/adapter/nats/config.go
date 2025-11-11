// Package nats
package nats

import (
	"time"
)

type Config struct {
	URL               string
	StreamName        string
	Subjects          []string
	MaxAge            time.Duration
	MaxBytes          int64
	Replicas          int
	ConsumerName      string
	DurableName       string
	AckWait           time.Duration
	MaxDeliver        int
	MaxAckPending     int
	ReplayPolicy      string
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration
}

func NewDefaultConfig() *Config {
	return &Config{
		URL:               "nats://localhost:4222",
		StreamName:        "EVENTS",
		Subjects:          []string{"events.>"}, // Aceita todos os eventos
		MaxAge:            7 * 24 * time.Hour,   // Mantém por 1 semana
		MaxBytes:          100 * 1024 * 1024,    // 100 MB
		Replicas:          1,                    // Sem replicação (dev)
		ConsumerName:      "default-consumer",
		DurableName:       "default-durable",
		AckWait:           30 * time.Second,
		MaxDeliver:        3,
		MaxAckPending:     100,
		ReplayPolicy:      "instant",
		ConnectionTimeout: 10 * time.Second,
		RequestTimeout:    5 * time.Second,
	}
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return ErrInvalidConfig("URL is required")
	}
	if c.StreamName == "" {
		return ErrInvalidConfig("StreamName is required")
	}
	if len(c.Subjects) == 0 {
		return ErrInvalidConfig("at least one Subject is required")
	}
	if c.MaxDeliver < 1 {
		return ErrInvalidConfig("MaxDeliver must be >= 1")
	}
	if c.MaxAckPending < 1 {
		return ErrInvalidConfig("MaxAckPending must be >= 1")
	}
	return nil
}
