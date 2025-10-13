// Package database
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewConnection(ctx context.Context, cfg *config.DatabaseConfig, logger repository.Logger) (*DB, error) {
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

	logger.Info(ctx, "database connection sucessfully")
	return &DB{Pool: pool}, nil
}

func (db *DB) Close(ctx context.Context, logger repository.Logger) {
	if db.Pool != nil {
		db.Pool.Close()
		logger.Info(ctx, "database connection closed")
	}
}
