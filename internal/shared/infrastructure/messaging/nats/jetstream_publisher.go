// Package nats
package nats

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type JetStreamPublisher struct {
	js jetstream.JetStream
}

func NewJetStreamPublisher(js jetstream.JetStream) *JetStreamPublisher {
	return &JetStreamPublisher{js: js}
}

func (p *JetStreamPublisher) Publish(
	ctx context.Context,
	subject string,
	IDMessage string,
	payload []byte,
) error {
	_, err := p.js.Publish(ctx, subject, payload, jetstream.WithMsgID(IDMessage))
	return err
}
