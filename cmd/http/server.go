// package main
package main

import (
	"context"
	"log"
	"log/slog"
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
	slog := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))
	logger := adapter.NewSlogAdapter(slog)
	db, err := postgres.NewConnection(ctx, &config.Database, logger)
	if err != nil {
		log.Fatal(err)
	}
	repository := persistence.NewAccount(db)
	// service := adapter.NewJWTAdapter(config.JWT)
	usecase := usecase.NewAuthenticate(repository, logger)
	handler := handler.NewSignIn(usecase, logger)

	mux := http.NewServeMux()
	// mux.Handle("POST /signin", http.HandlerFunc(handler.SignIn))
	mux.HandleFunc("POST /signin", handler.SignIn)
	// http.HandleFunc("POST /signin", handler.SignIn)
	http.ListenAndServe(":8080", middleware.LoggerMiddleware(logger)(mux))
}
