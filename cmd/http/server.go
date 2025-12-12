// package main
package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paladignus/actajus/internal/application/usecase"
	evt "github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/internal/infrastructure/adapter/messaging/nats"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/internal/infrastructure/database"
	"github.com/paladignus/actajus/internal/infrastructure/http/handler"
	"github.com/paladignus/actajus/internal/infrastructure/http/middleware"
	"github.com/paladignus/actajus/internal/infrastructure/messaging"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/infrastructure/security"
)

func main() {
	config := config.Load()
	ctx := context.Background()
	logger := adapter.NewDefaultLogger()
	db, err := database.NewConnection(ctx, &config.Database, logger)
	if err != nil {
		logger.Error(ctx, "error initializing the database connection.", "error", err)
		os.Exit(1)
	}
	defer db.Close(ctx, logger)
	persistence := postgres.NewPersistence(db)
	token := adapter.NewJWTAdapter(config.JWT)

	// CRIA EVENT REGISTRY E REGISTRA TIPOS DE EVENTOS
	// =================================================================
	registry := evt.NewRegistry()
	evt.RegisterAll(registry)

	// =================================================================
	// 2. CRIA PUBLISHER
	// =================================================================
	publisher, err := nats.NewPublisher(&config.NATS, logger)
	if err != nil {
		logger.Error(ctx, "❌ failed to create NATS publisher", "error", err)
		os.Exit(1)
	}
	logger.Info(ctx, "✅ NATS publisher connected")

	// =================================================================
	// 3. CRIA SUBSCRIBER
	// =================================================================
	subscriber, err := nats.NewSubscriber(&config.NATS, registry, logger)
	if err != nil {
		logger.Error(ctx, "❌ failed to create NATS subscriber", "error", err)
		os.Exit(1)
	}
	logger.Info(ctx, "✅ NATS subscriber connected")

	// =================================================================
	// 4. CRIA ADAPTERS (suas implementações existentes)
	// =================================================================
	smtp := adapter.NewSMTPEmail(config.SMTP)
	logger.Info(ctx, "✅ SMTP gateway configured")

	// =================================================================
	// 5. CRIA E REGISTRA HANDLERS
	// =================================================================
	// emailHandler := event.NewSendEmailHandler(&smtp)
	emailHandler := messaging.NewRecoverPassword(logger, smtp)
	// err = subscriber.Subscribe(
	// 	context.Background(), // Use context.Background() para Subscribe
	// 	"user.password_reset_requested",
	// 	emailHandler,
	// )
	if err != nil {
		logger.Error(ctx, "❌ failed to subscribe to events", "error", err)
		os.Exit(1)
	}
	logger.Info(ctx, "✅ email handler subscribed to password reset events")
	// =================================================================
	// 6. CRIA USECASES
	// =================================================================

	tokenService := security.NewHashToken()

	usecaseAuth := usecase.NewSignIn(
		persistence.Authentication(),
		logger,
		token,
	)
	usecaseGetEmail := usecase.NewGetEmailByCPF(persistence.Authentication(), logger)
	recoverPassword := usecase.NewRequestPasswordReset(persistence.Authentication(), persistence.Token(), tokenService, logger, publisher)
	sendEmailRecoverPassword := usecase.NewSendEmailRecoverPassword(
		logger,
		emailHandler,
		subscriber,
	)
	if err := sendEmailRecoverPassword.Execute(ctx); err != nil {
		logger.Error(ctx, "❌ failed to send email recover password", "error", err)
	}
	authHandler := handler.NewSignIn(usecaseAuth, logger)
	getEmailByCPF := handler.NewGetEmailByCPF(usecaseGetEmail, logger)
	recoverPasswordHandler := handler.NewRecoverPassword(recoverPassword, logger)
	renewPassword := handler.NewRenewPassword()

	// =================================================================
	// 7. CRIA E INICIA O SERVIDOR HTTP
	// =================================================================
	mux := http.NewServeMux()
	mux.HandleFunc("POST /signin", authHandler.SignIn)
	mux.HandleFunc("POST /auth/email", getEmailByCPF.GetEmailByCPF)
	mux.HandleFunc("GET /auth/recover/{token}", renewPassword.RenewPassword)
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
		logger.Info(ctx, "🚀 HTTP server running on port", "porta", config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "❌ failed to start the server", "error", err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "🛑 shutting down the server")
	if err := subscriber.Close(); err != nil {
		logger.Error(ctx, "error closing subscriber", "error", err)
	}
	if err := publisher.Close(); err != nil {
		logger.Error(ctx, "error closing publisher", "error", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, "error shutting down the server", "error", err)
		os.Exit(1)
	}
	logger.Info(ctx, "👋 shutdown complete")
}
