package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Bootstrap struct {
	Conn *nats.Conn
	JS   jetstream.JetStream
}

func NewBootstrap(url string) (*Bootstrap, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("create jetstream client: %w", err)
	}
	return &Bootstrap{
		Conn: nc,
		JS:   js,
	}, nil
}

func (b *Bootstrap) EnsureNotificationStream(ctx context.Context) error {
	_, err := b.JS.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      "NOTIFICATION",
		Subjects:  []string{"notification.>"},
		Retention: jetstream.WorkQueuePolicy,
		MaxAge:    24 * time.Hour,
		Storage:   jetstream.FileStorage,
		Replicas:  1,
	})
	if err != nil {
		return fmt.Errorf("create or update stream NOTIFICATION: %w", err)
	}
	return nil
}

func (b *Bootstrap) EnsureEmailWorkerConsumer(ctx context.Context) (jetstream.Consumer, error) {
	stream, err := b.JS.Stream(ctx, "NOTIFICATION")
	if err != nil {
		return nil, fmt.Errorf("get stream NOTIFICATION: %w", err)
	}
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          "email-worker",
		Durable:       "email-worker",
		FilterSubject: "notification.email.send_requested",
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       30 * time.Second,
		MaxDeliver:    5,
	})
	if err != nil {
		return nil, fmt.Errorf("create or update consumer email-worker: %w", err)
	}
	return consumer, nil
}
