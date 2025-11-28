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
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

type Publisher struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	config *config.NATSConfig
}

func NewPublisher(config *config.NATSConfig) (*Publisher, error) {
	opts := []nats.Option{
		nats.Timeout(config.ConnectionTimeout),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1), // Reconecta infinitamente
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				log.Printf("[NATS] Disconnected: %v\n", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("[NATS] Reconnected to %s\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("[NATS] Connection closed: %s\n", nc.Opts.Url)
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("[NATS] Server error: %v, subscription: %v\n", err, sub)
		}),
	}

	conn, err := nats.Connect(config.URL, opts...)
	if err != nil {
		return nil, ErrConnectionFailed(err)
	}
	js, err := conn.JetStream(nats.MaxWait(config.RequestTimeout))
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
		// Add more stream configuration options based on config
		Retention: nats.LimitsPolicy, // Default to limits policy
	}

	_, err := p.js.AddStream(streamConfig)
	if err != nil {
		// Check if it's the "stream name already in use" error
		if apiErr, ok := err.(*nats.APIError); ok && apiErr.ErrorCode == 10058 {
			// Stream already exists, try to update it
			_, err = p.js.UpdateStream(streamConfig)
			if err != nil {
				return ErrStreamCreationFailed(err)
			}
			return nil
		}
		return ErrStreamCreationFailed(err)
	}
	return nil
}

func (p *Publisher) Publish(ctx context.Context, evt event.IEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return ErrSerializationFailed(err)
	}
	
	subject := fmt.Sprintf("events.%s", evt.EventName())
	log.Printf("[NATS] Publishing event: %s (ID: %s)", subject, evt.AggregateID())
	
	// Use async publish with context timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.config.RequestTimeout)
	defer cancel()
	
	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  nats.Header{"Event-Type": []string{evt.EventName()}},
	}
	
	// Use async publish pattern for better performance
	pubAck, err := p.js.PublishMsgAsync(msg, nats.Context(ctxWithTimeout))
	if err != nil {
		return ErrPublishFailed(err)
	}
	
	// Optionally wait for the acknowledgment
	select {
	case <-pubAck.Ok():
		// Message acknowledged successfully
		return nil
	case err := <-pubAck.Err():
		return ErrPublishFailed(err)
	case <-ctxWithTimeout.Done():
		return ErrPublishFailed(fmt.Errorf("publish timeout"))
	}
}

func (p *Publisher) PublishBatch(ctx context.Context, events []event.IEvent) error {
	if len(events) == 0 {
		return nil
	}
	
	// For batch publishing, publish events individually
	var errors []error
	for _, evt := range events {
		if err := p.Publish(ctx, evt); err != nil {
			errors = append(errors, err)
			// Consider: should we continue or fail fast?
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("failed to publish %d/%d events: %v", len(errors), len(events), errors)
	}
	
	return nil
}

func (p *Publisher) Close() error {
	if p.conn != nil {
		// Wait for any async publishes to complete
		p.conn.Flush()
		p.conn.Close()
	}
	return nil
}

// HealthCheck returns the connection status
func (p *Publisher) HealthCheck(ctx context.Context) error {
	// Try a simple request to verify the connection is alive
	_, err := p.js.StreamInfo(p.config.StreamName, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	return nil
}