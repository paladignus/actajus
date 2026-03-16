// Package notification
package notification

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/paladignus/actajus/internal/module/notification/application/service"
	"github.com/paladignus/actajus/internal/module/notification/application/usecase"
	"github.com/paladignus/actajus/internal/module/notification/infrastructure/consumer"
	"github.com/paladignus/actajus/internal/module/notification/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/module/notification/infrastructure/template"
	sharedrepo "github.com/paladignus/actajus/internal/shared/application/repository"
	sharedsvc "github.com/paladignus/actajus/internal/shared/application/service"
	postgresShared "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type Dependencies struct {
	Logger        sharedrepo.Logger
	DB            postgresShared.Executor
	Serializer    sharedsvc.MessageSerializer
	Consumer      jetstream.Consumer
	EmailSender   service.EmailSender
	PublicBaseURL string
}

type Module struct {
	logger   sharedrepo.Logger
	consumer *consumer.EmailSendRequestedConsumer
}

func NewModule(dep Dependencies) (Module, error) {
	if dep.Consumer == nil {
		return Module{}, fmt.Errorf("notification: missing jetstream consumer")
	}
	if dep.EmailSender == nil {
		return Module{}, fmt.Errorf("notification: missing EmailSender")
	}
	if dep.PublicBaseURL == "" {
		return Module{}, fmt.Errorf("notification: missing PublicBaseURL")
	}
	sentRepo := postgres.NewSentEmailRepository(dep.DB)
	renderer, err := template.NewHTMLRenderer()
	if err != nil {
		return Module{}, fmt.Errorf("notification: %w", err)
	}
	uc := usecase.NewProcessEmailSendRequested(
		dep.EmailSender,
		renderer,
		sentRepo,
		dep.PublicBaseURL,
	)
	cons := consumer.NewEmailSendRequestedConsumer(
		dep.Consumer,
		dep.Serializer,
		uc,
		// dep.Logger,
	)
	return Module{
		logger:   dep.Logger,
		consumer: cons,
	}, nil
}

func (m Module) Start(ctx context.Context) {
	go func() {
		if err := m.consumer.Run(ctx); err != nil {
			m.logger.Error(context.Background(), "notification consumer stopped", "error", err)
		}
	}()
}
