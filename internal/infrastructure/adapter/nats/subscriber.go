// Package nats
package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/paladignus/actajus/internal/domain/event"
)

type Subscriber struct {
	conn          *nats.Conn
	js            nats.JetStreamContext
	config        *Config
	subscriptions map[string]*nats.Subscription // Mapeia eventName -> subscription
	mu            sync.RWMutex                  // Protege acesso ao map de subscriptions
	cancelFuncs   map[string]context.CancelFunc // Para cancelar goroutines
	wg            sync.WaitGroup                // Aguarda goroutines terminarem
	eventRegistry *event.EventRegistry
}

// var _ gateway.EventSubscriber = (*Subscriber)(nil)
func NewSubscriber(config *Config, registry *event.EventRegistry) (*Subscriber, error) { // Adicionado registry
	if err := config.Validate(); err != nil {
		return nil, err
	}
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
		subscriptions: make(map[string]*nats.Subscription),
		cancelFuncs:   make(map[string]context.CancelFunc),
		eventRegistry: registry,
	}, nil
}

func (p *Subscriber) Subscribe(ctx context.Context, eventName string, handler event.Handler) error {
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
		log.Printf("[NATS] Failed to create consumer: %s", consumerName)
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
	log.Printf("[NATS] Subscribed to event: %s (consumer: %s)", eventName, consumerName)
	return nil
}

func (p *Subscriber) processMessages(
	ctx context.Context,
	eventName string,
	sub *nats.Subscription,
	handler event.Handler,
) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Printf("[NATS] Stopping message processing for event: %s", eventName)
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

func (p *Subscriber) handleMessage(ctx context.Context, msg *nats.Msg, handler event.Handler) {
	processCtx, cancel := context.WithTimeout(ctx, p.config.AckWait-5*time.Second)
	defer cancel()
	meta, err := msg.Metadata()
	if err != nil {
		log.Printf("[NATS] Failed to get message metadata: %v", err)
		msg.Nak()
		return
	}
	log.Printf("[NATS] Processing message (attempt %d): %s", meta.NumDelivered, msg.Subject)
	var tempEvent map[string]interface{}
	if err := json.Unmarshal(msg.Data, &tempEvent); err != nil {
		log.Printf("[NATS] Failed to deserialize event header: %v", err)
		msg.Term()
		return
	}
	eventNameRaw, ok := tempEvent["Name"].(string)
	if !ok {
		log.Printf("[NATS] Failed to extract event name from message: %v", tempEvent)
		msg.Term()
		return
	}
	eventInstance, err := p.eventRegistry.CreateEventByName(eventNameRaw) // Usando p.eventRegistry
	if err != nil {
		log.Printf("[NATS] Failed to create event instance for name '%s': %v", eventNameRaw, err)
		msg.Term() // Descarta mensagem de evento desconhecido
		return
	}
	if err := json.Unmarshal(msg.Data, eventInstance); err != nil {
		log.Printf("[NATS] Failed to deserialize full event data for type '%s': %v", eventNameRaw, err)
		msg.Term()
		return
	}
	if !handler.CanHandle(eventInstance) {
		log.Printf("[NATS] Handler cannot process event: %s", eventInstance.EventName())
		msg.Ack()
		return
	}
	if err := handler.Handle(processCtx, eventInstance); err != nil {
		log.Printf("[NATS] Handler failed (attempt %d): %v", meta.NumDelivered, err)
		backoff := time.Duration(1<<(meta.NumDelivered-1)) * time.Second
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}
		log.Printf("[NATS] Retrying in %v...", backoff)
		msg.NakWithDelay(backoff)
		return
	}
	if err := msg.Ack(); err != nil {
		log.Printf("[NATS] Failed to ACK message: %v", err)
	} else {
		log.Printf("[NATS] Message processed successfully: %s", msg.Subject)
	}
}

func (p *Subscriber) Unsubscribe(eventName string) error {
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
	log.Printf("[NATS] Unsubscribed from event: %s", eventName)
	return nil
}

func (p *Subscriber) Close() error {
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
	log.Println("[NATS] Subscriber closed")
	return nil
}
