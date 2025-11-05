// Package main
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

func main() {
	config := config.Load()
	natsAdapter, err := adapter.NewNATSAdapter(config.NATS)
	if err != nil {
		log.Fatal("NATS connection failed:", err)
	}
	defer natsAdapter.Close()

	// smtp := &spy.SpySMTP{}
	// emailUC := usecase.NewSendWelcomeEmailUsecase(smtp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	smtp := adapter.NewSMTPEmail(config.SMTP)

	err = natsAdapter.SubscribeWithRetry(
		ctx,
		"user.created",
		"welcome-email-consumer", // nome durável
		func(msg []byte) error {
			return smtp.SendEmail(ctx, "marcelo@marcelo.eti.br", string(msg))
		},
	)
	// err = natsAdapter.SubscribeWithRetry(
	// 	ctx,
	// 	"user.created",
	// 	"welcome-email-consumer", // nome durável
	// 	func(msg []byte) error {
	// 		return emailUC.Handle(ctx, msg)
	// 	},
	// )
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}

	log.Println("📬 NATS consumer ready. Listening to 'user.created'...")

	// Esperar sinal para encerrar
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
}
