// package main
package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/paladignus/actajus/internal/application/usecase"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/internal/infrastructure/persistence"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
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
	adapter := adapter.NewJWTAdapter(config.JWT)
	signin := usecase.NewSignIn(
		repository,
		adapter,
	)
	output, err := signin.Execute(ctx, "727.753.511-15", "123456")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(output)
}
