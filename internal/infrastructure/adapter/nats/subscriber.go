// Package nats
package nats

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	js      nats.JetStreamContext
	dlqSubl string
}

type DLQMessage struct {
	OriginalSubject string    `json:"original_subject"`
	Data            []byte    `json:"data"`
	Error           string    `json:"error"`
	Attempts        uint64    `json:"attempts"`
	FailedAt        time.Time `json:"failed_at"`
}

func NewSubscriber(js nats.JetStreamContext, dlqSubj string) *Subscriber {
	return &Subscriber{js, dlqSubj}
}

// func (s *Subscriber) Subscribe(ctx context.Context, subject string, handler func([]byte) error) error {
// 	_, err := s.js.Subscribe(
// 		subject,
// 		func(m *nats.Msg) {
// 			log.Printf("📥 Message received: %s", string(m.Data))
// 			if err := handler(m.Data); err != nil {
// 				log.Printf("❌ Handler error: %v", err)
// 				m.NakWithDelay(5 * time.Second)
// 			} else {
// 				m.Ack()
// 			}
// 		},
// 		nats.ManualAck(),
// 		nats.AckExplicit(),
// 		nats.MaxDeliver(3),
// 		nats.AckWait(30*time.Second),
// 	)
// 	return err
// }

func (s *Subscriber) Subscribe(
	ctx context.Context,
	subject,
	durable string,
	handler func(msg []byte) error,
) error {
	_, err := s.js.Subscribe(
		subject,
		func(m *nats.Msg) {
			log.Printf("📥 Message received on subject '%s' (size: %d)", m.Subject, len(m.Data))
			meta, err := m.Metadata()
			if err != nil {
				log.Printf("failed to get message metadata: %v", err)
				m.Nak()
				return
			}
			log.Printf("📋 Delivery attempt: %d", meta.NumDelivered)
			if err := handler(m.Data); err != nil {
				log.Printf("❌ Handler failed: %v", err)
				attempts := meta.NumDelivered
				log.Printf("handler failed (attempt %d): failed %v", attempts, err)
				if attempts >= 3 {
					log.Printf("moving message to DLQ after %d attempts", attempts)
					s.sendToDLQ(m, err, attempts)
					m.Ack()
				} else {
					// Exponential backoff: 10s, 20s, 40s...
					delay := time.Duration(10<<attempts) * time.Second
					if delay > 60*time.Second {
						delay = 60 * time.Second
					}
					m.NakWithDelay(delay)
				}
			} else {
				log.Printf("✅ Handler succeeded")
				m.Ack() // ack na fila original para nao reprocessar
			}
		},
		nats.Durable(durable),
		nats.ManualAck(),
		nats.AckExplicit(),
		nats.MaxDeliver(10),
		nats.AckWait(20*time.Second),
	)
	return err
}

func (s *Subscriber) sendToDLQ(msg *nats.Msg, err error, attempts uint64) {
	dlq := DLQMessage{
		OriginalSubject: msg.Subject,
		Data:            msg.Data,
		Error:           err.Error(),
		Attempts:        attempts,
		FailedAt:        time.Now(),
	}
	payload, _ := json.Marshal(dlq)
	if _, err := s.js.Publish(s.dlqSubl, payload); err != nil {
		log.Printf("critical: failed to publish to DLQ: %v", err)
	}
}

func (s *Subscriber) ProcessDLQ(ctx context.Context, handler func(DLQMessage) error) error {
	_, err := s.js.Subscribe(
		s.dlqSubl,
		func(m *nats.Msg) {
			var dlq DLQMessage
			if err := json.Unmarshal(m.Data, &dlq); err != nil {
				log.Printf("invalid DLQ message: %v", err)
				m.Ack()
				return
			}
			if err := handler(dlq); err != nil {
				log.Printf("DLQ handler failed (will retry): %v", err)
				m.NakWithDelay(30 * time.Second)
				return
			}
			m.Ack()
		},
		nats.Durable("dlq-processor"),
		nats.ManualAck(),
	)
	return err
}
