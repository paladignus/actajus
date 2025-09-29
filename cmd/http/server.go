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
	"github.com/paladignus/actajus/internal/interface/controller"
	"github.com/paladignus/actajus/internal/interface/http/handler"
)

func main() {
	config := config.Load()
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := postgres.NewConnection(ctx, &config.Database, logger)
	if err != nil {
		log.Fatal(err)
	}
	repository := persistence.NewAuthenticate(db)
	service := adapter.NewJWTAdapter(config.JWT)
	usecase := usecase.NewSignIn(repository, service)
	controller := controller.NewSignIn(usecase)
	handler := handler.NewSignIn(controller)

	http.HandleFunc("POST /signin", handler.SignIn)
	http.ListenAndServe(":8080", nil)
}
