// Package nats
package nats

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

type Subscriber struct {
	conn          *nats.Conn
	js            nats.JetStreamContext
	config        *config.NATSConfig
	registry      *event.Registry               // Registry para deserializar eventos
	subscriptions map[string]*nats.Subscription // Mapeia eventName -> subscription
	mu            sync.RWMutex                  // Protege acesso ao map de subscriptions
	cancelFuncs   map[string]context.CancelFunc // Para cancelar goroutines
	wg            sync.WaitGroup                // Aguarda goroutines terminarem
	logger        repository.Logger
}

func NewSubscriber(config *config.NATSConfig, registry *event.Registry, logger repository.Logger) (*Subscriber, error) { // Adicionado registry
	conn, err := nats.Connect(
		config.URL,
		nats.Timeout(config.ConnectionTimeout),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, ErrConnectionFailed(err)
	}
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, ErrConnectionFailed(err)
	}
	return &Subscriber{
		conn:          conn,
		js:            js,
		config:        config,
		registry:      registry,
		subscriptions: make(map[string]*nats.Subscription),
		cancelFuncs:   make(map[string]context.CancelFunc),
		logger:        logger,
	}, nil
}

func (p *Subscriber) Subscribe(ctx context.Context, eventName string, handler repository.IHandler) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.subscriptions[eventName]; exists {
		return fmt.Errorf("already subscribed to event: %s", eventName)
	}
	subject := fmt.Sprintf("events.%s", eventName)
	consumerName := fmt.Sprintf("%s-%s", p.config.ConsumerName, strings.ReplaceAll(eventName, ".", "_"))
	consumerConfig := &nats.ConsumerConfig{
		Durable:       consumerName,
		DeliverPolicy: nats.DeliverAllPolicy,
		AckPolicy:     nats.AckExplicitPolicy,
		AckWait:       p.config.AckWait,
		MaxDeliver:    p.config.MaxDeliver,
		MaxAckPending: p.config.MaxAckPending,
		ReplayPolicy:  nats.ReplayInstantPolicy,
	}
	_, err := p.js.AddConsumer(p.config.StreamName, consumerConfig)
	if err != nil {
		p.logger.Error(ctx, "[NATS] Failed to create consumer", "consumer_name", consumerName, "error", err)
		return ErrConsumerCreationFailed(err)
	}
	sub, err := p.js.PullSubscribe(
		subject,
		consumerName,
		nats.ManualAck(),
	)
	if err != nil {
		return ErrSubscribeFailed(err)
	}
	p.subscriptions[eventName] = sub
	subCtx, cancel := context.WithCancel(ctx)
	p.cancelFuncs[eventName] = cancel
	p.wg.Add(1)
	go p.processMessages(subCtx, eventName, sub, handler)
	p.logger.Info(ctx, "[NATS] Subscribed to event", "event_name", eventName, "consumer", consumerName)
	return nil
}

func (p *Subscriber) processMessages(
	ctx context.Context,
	eventName string,
	sub *nats.Subscription,
	handler repository.IHandler,
) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			p.logger.Info(ctx, "[NATS] Stopping message processing for event", "event_name", eventName)
			return
		default:
			msgs, err := sub.Fetch(10, nats.Context(ctx))
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				time.Sleep(1 * time.Second)
				continue
			}
			for _, msg := range msgs {
				p.handleMessage(ctx, msg, handler)
			}
		}
	}
}

func (p *Subscriber) handleMessage(ctx context.Context, msg *nats.Msg, handler repository.IHandler) {
	meta, err := msg.Metadata()
	if err != nil {
		p.logger.Error(ctx, "[NATS] Failed to get message metadata", "error", err)
		msg.Nak() // Reenvia a mensagem
		return
	}
	p.logger.Info(ctx, "[NATS] Processing message", "num_delivered", meta.NumDelivered, "subject", msg.Subject)
	evt, err := p.registry.Unmarshal(msg.Data)
	if err != nil {
		p.logger.Error(ctx, "[NATS] Failed to deserialize event", "error", err)
		msg.Term()
		return
	}
	if err := handler.Handle(ctx, evt); err != nil {
		p.logger.Error(ctx, "[NATS] Handler failed", "num_delivered", meta.NumDelivered, "error", err)
		backoff := time.Duration(1<<(meta.NumDelivered-1)) * time.Second
		backoff = min(backoff, 30*time.Second)
		p.logger.Info(ctx, "[NATS] Retrying", "delay", backoff)
		msg.NakWithDelay(backoff)
		return
	}

	if err := msg.Ack(); err != nil {
		p.logger.Error(ctx, "[NATS] Failed to ACK message", "error", err)
	} else {
		p.logger.Info(ctx, "[NATS] Message processed successfully", "subject", msg.Subject)
	}
}

func (p *Subscriber) Unsubscribe(eventName string) error {
	ctx := context.Background()
	p.mu.Lock()
	defer p.mu.Unlock()
	sub, exists := p.subscriptions[eventName]
	if !exists {
		return fmt.Errorf("not subscribed to event: %s", eventName)
	}
	if cancel, ok := p.cancelFuncs[eventName]; ok {
		cancel()
		delete(p.cancelFuncs, eventName)
	}
	if err := sub.Unsubscribe(); err != nil {
		return err
	}
	delete(p.subscriptions, eventName)
	p.logger.Info(ctx, "[NATS] Unsubscribed from event", "event_name", eventName)
	return nil
}

func (p *Subscriber) Close() error {
	ctx := context.Background()
	p.mu.Lock()
	for _, cancel := range p.cancelFuncs {
		cancel()
	}
	for _, sub := range p.subscriptions {
		sub.Unsubscribe()
	}
	p.mu.Unlock()
	p.wg.Wait()
	if p.conn != nil {
		p.conn.Close()
	}
	p.logger.Info(ctx, "[NATS] Subscriber closed")
	return nil
}
