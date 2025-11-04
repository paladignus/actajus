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

	"github.com/paladignus/actajus/internal/application/usecase"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
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
	smtp := adapter.NewSMTPEmail(config.SMTP)
	usecaseAuth := usecase.NewSignIn(
		persistence.Authentication(),
		logger,
		token,
	)
	usecaseGetEmail := usecase.NewGetEmailByCPF(persistence.Authentication(), logger)

	recoverPassword := usecase.NewRecoverPassword(persistence.Authentication(), logger, token, &smtp)

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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}
	logger.Info(ctx, "✓ Servidor desligado com sucesso")
}
