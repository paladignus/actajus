// Package bootstrap provides shared application bootstrap for entrypoints.
package bootstrap

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/paladignus/actajus/internal/module/notification"
	"github.com/paladignus/actajus/internal/module/notification/infrastructure/email"
	sharednats "github.com/paladignus/actajus/internal/shared/infrastructure/messaging/nats"
	sharedoutbox "github.com/paladignus/actajus/internal/shared/infrastructure/messaging/outbox"
)

// NotificationRuntime holds the messaging infrastructure required by the notification module.
type NotificationRuntime struct {
	Base *Runtime
	NATS *sharednats.Bootstrap
}

// NewNotificationRuntime initializes NATS, outbox dispatching and the notification consumer.
func NewNotificationRuntime(ctx context.Context, base *Runtime) (*NotificationRuntime, error) {
	natsBoot, err := sharednats.NewBootstrap(base.Config.NATS.URL)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}

	if err := natsBoot.EnsureNotificationStream(ctx); err != nil {
		natsBoot.Conn.Close()
		return nil, fmt.Errorf("ensure notification stream: %w", err)
	}

	emailConsumer, err := natsBoot.EnsureEmailWorkerConsumer(ctx)
	if err != nil {
		natsBoot.Conn.Close()
		return nil, fmt.Errorf("ensure email consumer: %w", err)
	}

	publisher := sharednats.NewJetStreamPublisher(natsBoot.JS)
	dispatcher := sharedoutbox.NewDispatcher(base.DB, publisher, 200)
	go func() {
		if err := dispatcher.Run(ctx, 2*time.Second); err != nil && err != context.Canceled {
			base.Logger.Error(context.Background(), "outbox dispatcher stopped", "error", err)
		}
	}()

	smtpPort, err := strconv.Atoi(base.Config.SMTP.Port)
	if err != nil {
		natsBoot.Conn.Close()
		return nil, fmt.Errorf("invalid SMTP_PORT %q: %w", base.Config.SMTP.Port, err)
	}

	emailSender := email.NewSMTPSender(
		base.Config.SMTP.Host,
		smtpPort,
		base.Config.SMTP.User,
		base.Config.SMTP.Pass,
		base.Config.SMTP.From,
	)

	module, err := notification.NewModule(notification.Dependencies{
		Logger:        base.Logger,
		DB:            base.DB,
		Serializer:    base.Codec,
		Consumer:      emailConsumer,
		EmailSender:   emailSender,
		PublicBaseURL: base.Config.SMTP.PublicBaseURL,
	})
	if err != nil {
		natsBoot.Conn.Close()
		return nil, fmt.Errorf("init notification module: %w", err)
	}
	module.Start(ctx)

	return &NotificationRuntime{
		Base: base,
		NATS: natsBoot,
	}, nil
}

// Close releases the messaging resources owned by the notification runtime.
func (r *NotificationRuntime) Close() {
	if r.NATS != nil && r.NATS.Conn != nil {
		r.NATS.Conn.Close()
		r.NATS = nil
	}
}
