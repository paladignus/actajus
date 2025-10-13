// package main
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paladignus/actajus/internal/application/usecase"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/internal/infrastructure/controller"
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
	usecase := usecase.NewAuthenticate(
		persistence.Account(),
		logger,
		token,
	)
	service := controller.NewAuthenticate(usecase)
	authHandler := handler.NewSignIn(service, logger)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /signin", authHandler.SignIn)
	handler := middleware.LoggerMiddleware(logger)(mux)
	srv := &http.Server{
		Addr:         ":" + config.Server.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
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
