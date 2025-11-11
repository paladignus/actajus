// Package nats
package nats

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/gateway"
)

// Subscriber implementa gateway.EventSubscriber usando NATS JetStream.
// Processa eventos com retry automático e dead letter queue.
type Subscriber struct {
	conn          *nats.Conn
	js            nats.JetStreamContext
	config        *Config
	registry      *event.Registry               // Registry para deserializar eventos
	subscriptions map[string]*nats.Subscription // Mapeia eventName -> subscription
	mu            sync.RWMutex                  // Protege acesso ao map de subscriptions
	cancelFuncs   map[string]context.CancelFunc // Para cancelar goroutines
	wg            sync.WaitGroup                // Aguarda goroutines terminarem
}

var _ gateway.EventSubscriber = (*Subscriber)(nil)

// NewSubscriber cria um novo Subscriber conectado ao NATS.
// Opcionalmente aceita um registry customizado, caso contrário usa o global.
func NewSubscriber(config *Config, registry ...*event.Registry) (*Subscriber, error) {
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

	// Usa registry customizado ou global
	var reg *event.Registry
	if len(registry) > 0 {
		reg = registry[0]
	} else {
		reg = event.GlobalRegistry
	}

	return &Subscriber{
		conn:          conn,
		js:            js,
		config:        config,
		registry:      reg,
		subscriptions: make(map[string]*nats.Subscription),
		cancelFuncs:   make(map[string]context.CancelFunc),
	}, nil
}

// Subscribe cria um consumer durável e começa a processar eventos.
//
// Parâmetros:
//   - ctx: contexto para cancelamento (use context.Background() normalmente)
//   - eventName: nome do evento (ex: "user.password_reset_requested")
//   - handler: handler que processará os eventos
//
// Comportamento:
//  1. Cria um consumer durável no JetStream
//  2. Inicia goroutine processando mensagens continuamente
//  3. Para cada mensagem:
//     - Desserializa o evento
//     - Chama o handler
//     - Se sucesso: envia ACK
//     - Se falha: envia NAK (retry automático)
//  4. Após MaxDeliver tentativas, mensagem vai para DLQ
func (p *Subscriber) Subscribe(ctx context.Context, eventName string, handler event.Handler) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Verifica se já existe subscription para este evento
	if _, exists := p.subscriptions[eventName]; exists {
		return fmt.Errorf("already subscribed to event: %s", eventName)
	}

	// Cria subject específico para este evento
	subject := fmt.Sprintf("events.%s", eventName)

	// Cria consumer durável
	// consumerName := fmt.Sprintf("%s-%s", p.config.ConsumerName, eventName)
	consumerName := fmt.Sprintf("%s-%s", p.config.ConsumerName, strings.ReplaceAll(eventName, ".", "_"))

	consumerConfig := &nats.ConsumerConfig{
		Durable:       consumerName,
		DeliverPolicy: nats.DeliverAllPolicy,  // Entrega todas as mensagens
		AckPolicy:     nats.AckExplicitPolicy, // Requer ACK explícito
		AckWait:       p.config.AckWait,
		MaxDeliver:    p.config.MaxDeliver,
		MaxAckPending: p.config.MaxAckPending,
		ReplayPolicy:  nats.ReplayInstantPolicy,
	}

	// Cria ou atualiza o consumer
	_, err := p.js.AddConsumer(p.config.StreamName, consumerConfig)
	if err != nil {
		return ErrConsumerCreationFailed(err)
	}

	// Cria subscription pull-based (mais controle sobre processamento)
	sub, err := p.js.PullSubscribe(
		subject,
		consumerName,
		nats.ManualAck(), // ACK manual para controle fino
	)
	if err != nil {
		return ErrSubscribeFailed(err)
	}

	p.subscriptions[eventName] = sub

	// Cria contexto cancelável para esta subscription
	subCtx, cancel := context.WithCancel(ctx)
	p.cancelFuncs[eventName] = cancel

	// Inicia goroutine para processar mensagens
	p.wg.Add(1)
	go p.processMessages(subCtx, eventName, sub, handler)

	log.Printf("[NATS] Subscribed to event: %s (consumer: %s)", eventName, consumerName)
	return nil
}

// processMessages processa mensagens continuamente em background.
// Esta função roda em uma goroutine separada.
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
			// Puxa até 10 mensagens por vez (batch processing)
			msgs, err := sub.Fetch(10, nats.Context(ctx))
			if err != nil {
				// Se contexto cancelado, sai do loop
				if ctx.Err() != nil {
					return
				}
				// Outros erros: aguarda antes de tentar novamente
				time.Sleep(1 * time.Second)
				continue
			}

			// Processa cada mensagem
			for _, msg := range msgs {
				p.handleMessage(ctx, msg, handler)
			}
		}
	}
}

// handleMessage processa uma única mensagem com retry inteligente.
func (p *Subscriber) handleMessage(ctx context.Context, msg *nats.Msg, handler event.Handler) {
	// Cria contexto com timeout para processamento
	processCtx, cancel := context.WithTimeout(ctx, p.config.AckWait-5*time.Second)
	defer cancel()

	// Obtém metadata da mensagem para logging
	meta, err := msg.Metadata()
	if err != nil {
		log.Printf("[NATS] Failed to get message metadata: %v", err)
		msg.Nak() // Reenvia a mensagem
		return
	}

	// Log de tentativa
	log.Printf("[NATS] Processing message (attempt %d): %s",
		meta.NumDelivered, msg.Subject)

	// Desserializa o evento usando o Registry (TYPE-SAFE!)
	evt, err := p.registry.Unmarshal(msg.Data)
	if err != nil {
		log.Printf("[NATS] Failed to deserialize event: %v", err)
		// Mensagem inválida ou tipo não registrado - envia para DLQ (término)
		msg.Term()
		return
	}

	// Verifica se o handler pode processar este evento
	if !handler.CanHandle(evt) {
		log.Printf("[NATS] Handler cannot process event: %s", evt.EventName())
		msg.Ack() // ACK para não reprocessar
		return
	}

	// Executa o handler
	if err := handler.Handle(processCtx, evt); err != nil {
		log.Printf("[NATS] Handler failed (attempt %d): %v",
			meta.NumDelivered, err)

		// Calcula backoff exponencial para retry
		// Tentativa 1: 1s, Tentativa 2: 2s, Tentativa 3: 4s, etc
		backoff := time.Duration(1<<(meta.NumDelivered-1)) * time.Second
		if backoff > 30*time.Second {
			backoff = 30 * time.Second // Máximo de 30s
		}

		log.Printf("[NATS] Retrying in %v...", backoff)

		// NAK com delay (retry com backoff)
		msg.NakWithDelay(backoff)
		return
	}

	// Sucesso - envia ACK
	if err := msg.Ack(); err != nil {
		log.Printf("[NATS] Failed to ACK message: %v", err)
	} else {
		log.Printf("[NATS] Message processed successfully: %s", msg.Subject)
	}
}

// Unsubscribe cancela a inscrição em um evento específico.
func (p *Subscriber) Unsubscribe(eventName string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	sub, exists := p.subscriptions[eventName]
	if !exists {
		return fmt.Errorf("not subscribed to event: %s", eventName)
	}

	// Cancela a goroutine de processamento
	if cancel, ok := p.cancelFuncs[eventName]; ok {
		cancel()
		delete(p.cancelFuncs, eventName)
	}

	// Unsubscribe
	if err := sub.Unsubscribe(); err != nil {
		return err
	}

	delete(p.subscriptions, eventName)
	log.Printf("[NATS] Unsubscribed from event: %s", eventName)
	return nil
}

// Close fecha todas as subscriptions e a conexão.
func (p *Subscriber) Close() error {
	p.mu.Lock()

	// Cancela todas as goroutines
	for _, cancel := range p.cancelFuncs {
		cancel()
	}

	// Unsubscribe de tudo
	for _, sub := range p.subscriptions {
		sub.Unsubscribe()
	}

	p.mu.Unlock()

	// Aguarda todas as goroutines terminarem
	p.wg.Wait()

	// Fecha conexão
	if p.conn != nil {
		p.conn.Close()
	}

	log.Println("[NATS] Subscriber closed")
	return nil
}

// type Subscriber struct {
// 	conn          *nats.Conn
// 	js            nats.JetStreamContext
// 	config        *Config
// 	subscriptions map[string]*nats.Subscription // Mapeia eventName -> subscription
// 	mu            sync.RWMutex                  // Protege acesso ao map de subscriptions
// 	cancelFuncs   map[string]context.CancelFunc // Para cancelar goroutines
// 	wg            sync.WaitGroup                // Aguarda goroutines terminarem
// 	eventRegistry *event.EventRegistry
// }
//
// // var _ gateway.EventSubscriber = (*Subscriber)(nil)
// func NewSubscriber(config *Config, registry *event.EventRegistry) (*Subscriber, error) { // Adicionado registry
// 	if err := config.Validate(); err != nil {
// 		return nil, err
// 	}
// 	conn, err := nats.Connect(
// 		config.URL,
// 		nats.Timeout(config.ConnectionTimeout),
// 		nats.RetryOnFailedConnect(true),
// 		nats.MaxReconnects(-1),
// 		nats.ReconnectWait(2*time.Second),
// 	)
// 	if err != nil {
// 		return nil, ErrConnectionFailed(err)
// 	}
// 	js, err := conn.JetStream()
// 	if err != nil {
// 		conn.Close()
// 		return nil, ErrConnectionFailed(err)
// 	}
// 	return &Subscriber{
// 		conn:          conn,
// 		js:            js,
// 		config:        config,
// 		subscriptions: make(map[string]*nats.Subscription),
// 		cancelFuncs:   make(map[string]context.CancelFunc),
// 		eventRegistry: registry,
// 	}, nil
// }
//
// func (p *Subscriber) Subscribe(ctx context.Context, eventName string, handler event.Handler) error {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	if _, exists := p.subscriptions[eventName]; exists {
// 		return fmt.Errorf("already subscribed to event: %s", eventName)
// 	}
// 	subject := fmt.Sprintf("events.%s", eventName)
// 	consumerName := fmt.Sprintf("%s-%s", p.config.ConsumerName, strings.ReplaceAll(eventName, ".", "_"))
// 	consumerConfig := &nats.ConsumerConfig{
// 		Durable:       consumerName,
// 		DeliverPolicy: nats.DeliverAllPolicy,
// 		AckPolicy:     nats.AckExplicitPolicy,
// 		AckWait:       p.config.AckWait,
// 		MaxDeliver:    p.config.MaxDeliver,
// 		MaxAckPending: p.config.MaxAckPending,
// 		ReplayPolicy:  nats.ReplayInstantPolicy,
// 	}
// 	_, err := p.js.AddConsumer(p.config.StreamName, consumerConfig)
// 	if err != nil {
// 		log.Printf("[NATS] Failed to create consumer: %s", consumerName)
// 		return ErrConsumerCreationFailed(err)
// 	}
// 	sub, err := p.js.PullSubscribe(
// 		subject,
// 		consumerName,
// 		nats.ManualAck(),
// 	)
// 	if err != nil {
// 		return ErrSubscribeFailed(err)
// 	}
// 	p.subscriptions[eventName] = sub
// 	subCtx, cancel := context.WithCancel(ctx)
// 	p.cancelFuncs[eventName] = cancel
// 	p.wg.Add(1)
// 	go p.processMessages(subCtx, eventName, sub, handler)
// 	log.Printf("[NATS] Subscribed to event: %s (consumer: %s)", eventName, consumerName)
// 	return nil
// }
//
// func (p *Subscriber) processMessages(
// 	ctx context.Context,
// 	eventName string,
// 	sub *nats.Subscription,
// 	handler event.Handler,
// ) {
// 	defer p.wg.Done()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			log.Printf("[NATS] Stopping message processing for event: %s", eventName)
// 			return
// 		default:
// 			msgs, err := sub.Fetch(10, nats.Context(ctx))
// 			if err != nil {
// 				if ctx.Err() != nil {
// 					return
// 				}
// 				time.Sleep(1 * time.Second)
// 				continue
// 			}
// 			for _, msg := range msgs {
// 				p.handleMessage(ctx, msg, handler)
// 			}
// 		}
// 	}
// }
//
// func (p *Subscriber) handleMessage(ctx context.Context, msg *nats.Msg, handler event.Handler) {
// 	processCtx, cancel := context.WithTimeout(ctx, p.config.AckWait-5*time.Second)
// 	defer cancel()
// 	meta, err := msg.Metadata()
// 	if err != nil {
// 		log.Printf("[NATS] Failed to get message metadata: %v", err)
// 		msg.Nak()
// 		return
// 	}
// 	log.Printf("[NATS] Processing message (attempt %d): %s", meta.NumDelivered, msg.Subject)
// 	var tempEvent map[string]interface{}
// 	if err := json.Unmarshal(msg.Data, &tempEvent); err != nil {
// 		log.Printf("[NATS] Failed to deserialize event header: %v", err)
// 		msg.Term()
// 		return
// 	}
// 	eventNameRaw, ok := tempEvent["Name"].(string)
// 	if !ok {
// 		log.Printf("[NATS] Failed to extract event name from message: %v", tempEvent)
// 		msg.Term()
// 		return
// 	}
// 	eventInstance, err := p.eventRegistry.CreateEventByName(eventNameRaw) // Usando p.eventRegistry
// 	if err != nil {
// 		log.Printf("[NATS] Failed to create event instance for name '%s': %v", eventNameRaw, err)
// 		msg.Term() // Descarta mensagem de evento desconhecido
// 		return
// 	}
// 	if err := json.Unmarshal(msg.Data, eventInstance); err != nil {
// 		log.Printf("[NATS] Failed to deserialize full event data for type '%s': %v", eventNameRaw, err)
// 		msg.Term()
// 		return
// 	}
// 	if !handler.CanHandle(eventInstance) {
// 		log.Printf("[NATS] Handler cannot process event: %s", eventInstance.EventName())
// 		msg.Ack()
// 		return
// 	}
// 	if err := handler.Handle(processCtx, eventInstance); err != nil {
// 		log.Printf("[NATS] Handler failed (attempt %d): %v", meta.NumDelivered, err)
// 		backoff := time.Duration(1<<(meta.NumDelivered-1)) * time.Second
// 		if backoff > 30*time.Second {
// 			backoff = 30 * time.Second
// 		}
// 		log.Printf("[NATS] Retrying in %v...", backoff)
// 		msg.NakWithDelay(backoff)
// 		return
// 	}
// 	if err := msg.Ack(); err != nil {
// 		log.Printf("[NATS] Failed to ACK message: %v", err)
// 	} else {
// 		log.Printf("[NATS] Message processed successfully: %s", msg.Subject)
// 	}
// }
//
// func (p *Subscriber) Unsubscribe(eventName string) error {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	sub, exists := p.subscriptions[eventName]
// 	if !exists {
// 		return fmt.Errorf("not subscribed to event: %s", eventName)
// 	}
// 	if cancel, ok := p.cancelFuncs[eventName]; ok {
// 		cancel()
// 		delete(p.cancelFuncs, eventName)
// 	}
// 	if err := sub.Unsubscribe(); err != nil {
// 		return err
// 	}
// 	delete(p.subscriptions, eventName)
// 	log.Printf("[NATS] Unsubscribed from event: %s", eventName)
// 	return nil
// }
//
// func (p *Subscriber) Close() error {
// 	p.mu.Lock()
// 	for _, cancel := range p.cancelFuncs {
// 		cancel()
// 	}
// 	for _, sub := range p.subscriptions {
// 		sub.Unsubscribe()
// 	}
// 	p.mu.Unlock()
// 	p.wg.Wait()
// 	if p.conn != nil {
// 		p.conn.Close()
// 	}
// 	log.Println("[NATS] Subscriber closed")
// 	return nil
// }
