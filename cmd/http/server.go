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
	"github.com/paladignus/actajus/internal/infrastructure/http/handler"
	"github.com/paladignus/actajus/internal/infrastructure/http/middleware"
	"github.com/paladignus/actajus/internal/infrastructure/messaging"
	"github.com/paladignus/actajus/internal/infrastructure/metrics"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/infrastructure/security"
	"github.com/paladignus/actajus/internal/infrastructure/tracing"
)

func main() {
	config := config.Load()
	ctx := context.Background()
	logger := adapter.NewDefaultLogger()

	// Initialize distributed tracing if enabled
	if config.Tracing.Enabled {
		logger.Info(ctx, "Initializing distributed tracing",
			"service_name", config.Tracing.ServiceName,
			"agent_host", config.Tracing.AgentHost,
			"agent_port", config.Tracing.AgentPort)

		if err := tracing.InitTracer(
			config.Tracing.ServiceName,
			config.Tracing.AgentHost,
			config.Tracing.AgentPort,
		); err != nil {
			logger.Error(ctx, "Failed to initialize tracing", "error", err)
			// Continue execution even if tracing fails to initialize
		} else {
			logger.Info(ctx, "Distributed tracing initialized successfully")

			// Ensure tracing is properly shut down on exit
			defer func() {
				if err := tracing.ShutdownTracer(ctx); err != nil {
					logger.Error(ctx, "Failed to shutdown tracer", "error", err)
				} else {
					logger.Info(ctx, "Tracer shutdown successfully")
				}
			}()
		}
	}

	db, err := postgres.NewConnection(ctx, &config.Database, logger)
	if err != nil {
		logger.Error(ctx, "error initializing the database connection.", "error", err)
		os.Exit(1)
	}
	// defer db.Close(ctx, logger)
	defer db.Close()
	// persistence := postgres.NewPersistence(db)
	// uow := postgres.NewUnitOfWork(db)
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

	tokenService := security.NewCryptoTokenGenerator()
	// passwordService := security.NewBcryptPasswordHasher(10)

	usecaseAuth := usecase.NewSignIn(
		persistence.User(),
		token,
	)
	usecaseGetEmail := usecase.NewGetEmailByCPF(persistence.User())
	recoverPassword := usecase.NewRequestPasswordReset(
		persistence.User(),
		persistence.PasswordResetToken(),
		tokenService,
		publisher,
	)
	sendEmailRecoverPassword := usecase.NewSendEmailRecoverPassword(
		logger,
		emailHandler,
		subscriber,
	)
	if err := sendEmailRecoverPassword.Execute(ctx); err != nil {
		logger.Error(ctx, "❌ failed to send email recover password", "error", err)
	}

	renewPasswordUC := usecase.NewRenewPassword(
		persistence.User(),
		persistence.PasswordResetToken(),
	)

	u := postgres.NewUnitOfWorkEnterprise(db)
	enterpriseCreateUC := usecase.NewEnterprise(*u)

	// =================================================================
	// 7. CRIA OS HANDLERS
	// =================================================================

	authHandler := handler.NewSignIn(usecaseAuth, logger)
	getEmailByCPF := handler.NewGetEmailByCPF(usecaseGetEmail, logger)
	resetPasswordHandler := handler.NewRequestPasswordReset(recoverPassword, logger)
	renewPasswordHandler := handler.NewRenewPassword(renewPasswordUC, logger)
	enterpriseHandler := handler.NewEnterprise(enterpriseCreateUC, logger)

	// =================================================================
	// 8. CRIA E INICIA O SERVIDOR HTTP
	// =================================================================
	mux := http.NewServeMux()
	mux.HandleFunc("POST /signin", authHandler.SignIn)
	mux.HandleFunc("POST /auth/email", getEmailByCPF.GetEmailByCPF)
	mux.HandleFunc("POST /auth/recover", resetPasswordHandler.RequestPasswordReset)
	mux.HandleFunc("POST /auth/renew", renewPasswordHandler.RenewPassword)
	mux.HandleFunc("POST /enterprise", enterpriseHandler.Create)

	// Add metrics endpoint for Prometheus
	mux.Handle("/metrics", metrics.Handler())

	// Create handler chain with multiple middleware:
	// 1. Metrics middleware to collect Prometheus metrics
	// 2. Tracing middleware to capture distributed traces
	// 3. Logging middleware to capture request logs
	// 4. CORS middleware to handle cross-origin requests
	handler := middleware.EnableCORS(
		middleware.LoggerMiddleware(logger)(
			tracing.Middleware(config.Tracing.ServiceName)(
				metrics.Middleware(mux),
			),
		),
	)

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
