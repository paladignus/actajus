// package main
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/paladignus/actajus/internal/application/usecase"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/internal/infrastructure/persistence"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/interface/http/handler"
	"github.com/paladignus/actajus/internal/interface/http/middleware"
)

func main() {
	config := config.Load()
	ctx := context.Background()
	logger := adapter.NewDefaultLogger()
	db, err := postgres.NewConnection(ctx, &config.Database, logger)
	if err != nil {
		log.Fatal(err)
	}
	persistence := persistence.NewPersistence(db)
	token := adapter.NewJWTAdapter(config.JWT)
	usecase := usecase.NewAuthenticate(
		persistence.Account(),
		logger,
		token,
	)
	authHandler := handler.NewSignIn(usecase, logger)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /signin", authHandler.SignIn)
	logger.Info(ctx, "starting server", "port", config.Server.Port)
	if err := http.ListenAndServe(":"+config.Server.Port, middleware.LoggerMiddleware(logger)(mux)); err != nil {
		logger.Error(ctx, "server failed", "error", err)
		os.Exit(1)
	}
}
