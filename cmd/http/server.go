// package main
package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paladignus/actajus/internal/application/event"
	"github.com/paladignus/actajus/internal/application/usecase"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/internal/infrastructure/adapter/nats"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/internal/infrastructure/database"
	"github.com/paladignus/actajus/internal/infrastructure/http/handler"
	"github.com/paladignus/actajus/internal/infrastructure/http/middleware"
	"github.com/paladignus/actajus/internal/infrastructure/persistence"
)

func main() {
	config := config.Load()
	ctx := context.Background()
	logger := adapter.NewDefaultLogger()
	db, err := database.NewConnection(ctx, &config.Database, logger)
	if err != nil {
		log.Fatalf("erro ao inicializar a aplicação: %v", err)
	}
	defer db.Close(ctx, logger)
	persistence := persistence.NewPersistence(db)
	token := adapter.NewJWTAdapter(config.JWT)

	natsConfig := &nats.Config{
		URL:               "nats://localhost:4222",
		StreamName:        "EVENTS",
		Subjects:          []string{"events.>"}, // Aceita todos os eventos
		MaxAge:            7 * 24 * time.Hour,   // Mantém por 1 semana
		MaxBytes:          100 * 1024 * 1024,    // 100 MB
		Replicas:          1,
		ConsumerName:      "actajus-consumer",
		DurableName:       "actajus-durable",
		AckWait:           30 * time.Second,
		MaxDeliver:        3, // Tenta até 3 vezes
		MaxAckPending:     100,
		ReplayPolicy:      "instant",
		ConnectionTimeout: 10 * time.Second,
		RequestTimeout:    5 * time.Second,
	}
	// =================================================================
	// CRIA EVENT REGISTRY E REGISTRA TIPOS DE EVENTOS
	// =================================================================
	// eventRegistry := evt.NewEventRegistry()
	// Registra os tipos de eventos conhecidos explicitamente aqui.
	// Isso substitui o uso de init() e torna o processo explícito.
	// err = eventRegistry.Register("user.password_reset_requested", func() evt.Event { return &evt.PasswordResetRequestedEvent{} })
	// if err != nil {
	// 	log.Fatalf("❌ Failed to register event type: %v", err)
	// }
	// Futuramente, ao adicionar um novo evento:
	// err = eventRegistry.Register("user.created", func() event.Event { return &event.UserCreatedEvent{} })
	// if err != nil { ... }

	// Log de eventos registrados (opcional)
	// log.Println("Registered Events:", eventRegistry.GetRegisteredEventNames())
	// =================================================================

	// =================================================================
	// 2. CRIA PUBLISHER
	// =================================================================
	publisher, err := nats.NewPublisher(natsConfig)
	if err != nil {
		log.Fatalf("❌ Failed to create NATS publisher: %v", err)
	}
	defer publisher.Close()
	log.Println("✅ NATS Publisher connected")

	// =================================================================
	// 3. CRIA SUBSCRIBER
	// =================================================================
	subscriber, err := nats.NewSubscriber(natsConfig)
	if err != nil {
		log.Fatalf("❌ Failed to create NATS subscriber: %v", err)
	}
	defer subscriber.Close()
	log.Println("✅ NATS Subscriber connected")

	// =================================================================
	// 4. CRIA ADAPTERS (suas implementações existentes)
	// =================================================================
	smtpGateway := adapter.NewSMTPEmail(config.SMTP)
	log.Println("✅ SMTP Gateway configured")

	// =================================================================
	// 5. CRIA E REGISTRA HANDLERS
	// =================================================================
	emailHandler := event.NewSendEmailHandler(&smtpGateway)

	err = subscriber.Subscribe(
		context.Background(), // Use context.Background() para Subscribe
		"user.password_reset_requested",
		emailHandler,
	)
	if err != nil {
		log.Fatalf("❌ Failed to subscribe to events: %v", err)
	}
	log.Println("✅ Email handler subscribed to password reset events")

	usecaseAuth := usecase.NewSignIn(
		persistence.Authentication(),
		logger,
		token,
	)
	usecaseGetEmail := usecase.NewGetEmailByCPF(persistence.Authentication(), logger)

	recoverPassword := usecase.NewRecoverPassword(persistence.Authentication(), logger, token, publisher)

	authHandler := handler.NewSignIn(usecaseAuth, logger)

	getEmailByCPF := handler.NewGetEmailByCPF(usecaseGetEmail, logger)

	recoverPasswordHandler := handler.NewRecoverPassword(recoverPassword, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signin", authHandler.SignIn)
	mux.HandleFunc("POST /auth/email", getEmailByCPF.GetEmailByCPF)
	mux.HandleFunc("POST /auth/recover", recoverPasswordHandler.RecoverPassword)

	handler := middleware.EnableCORS(middleware.LoggerMiddleware(logger)(mux))
	srv := &http.Server{
		Addr:         ":" + config.Server.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	go func() {
		logger.Info(ctx, "🚀 Servidor HTTP rodando na porta", "porta", config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "🛑 Desligando servidor...")
	// Fecha conexões na ordem correta
	if err := subscriber.Close(); err != nil {
		log.Printf("Error closing subscriber: %v", err)
	}
	if err := publisher.Close(); err != nil {
		log.Printf("Error closing publisher: %v", err)
	}

	log.Println("👋 Shutdown complete")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}
	logger.Info(ctx, "✓ Servidor desligado com sucesso")
}
