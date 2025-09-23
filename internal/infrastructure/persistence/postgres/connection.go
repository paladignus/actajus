// Package postgres
package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

type DB struct {
	Pool *pgxpool.Pool
	log  *slog.Logger
}

func NewConnection(ctx context.Context, cfg *config.DatabaseConfig, logger *slog.Logger) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	logger.Info("database connection sucessfully")
	return &DB{Pool: pool, log: logger}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
	db.log.Info("database connection closed")
}
