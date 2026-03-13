// Package consumer
package consumer

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/paladignus/actajus/internal/module/notification/application/dto"
	"github.com/paladignus/actajus/internal/module/notification/application/usecase"
	"github.com/paladignus/actajus/internal/shared/application/service"
)

type EmailSendRequestedConsumer struct {
	consumer   jetstream.Consumer
	serializer service.MessageSerializer
	usecase    usecase.ProcessEmailSendRequested
}

func NewEmailSendRequestedConsumer(
	consumer jetstream.Consumer,
	serializer service.MessageSerializer,
	usecase usecase.ProcessEmailSendRequested,
) *EmailSendRequestedConsumer {
	return &EmailSendRequestedConsumer{
		consumer,
		serializer,
		usecase,
	}
}

func (c *EmailSendRequestedConsumer) Run(ctx context.Context) error {
	msgs, err := c.consumer.Messages()
	if err != nil {
		return fmt.Errorf("consumer messages: %w", err)
	}
	defer msgs.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := msgs.Next()
			if err != nil {
				return fmt.Errorf("next message: %w", err)
			}
			var payload dto.EmailSendRequested
			if err := c.serializer.Unmarshal(msg.Data(), &payload); err != nil {
				_ = msg.Term()
				continue
			}
			if err := c.usecase.Execute(ctx, payload); err != nil {
				_ = msg.Nak()
				continue
			}
			if err := msg.Ack(); err != nil {
				return fmt.Errorf("ack message: %w", err)
			}
		}
	}
}
