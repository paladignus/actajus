// Package grpcserver provides the gRPC/Connect server
package grpcserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/paladignus/actajus/pkg/di"
)

// Server represents the gRPC server
type Server struct {
	container *di.Container
	srv       *http.Server
}

// New creates a new gRPC server
func New(container *di.Container) *Server {
	mux := http.NewServeMux()

	handler := corsMiddleware(mux)

	return &Server{
		container: container,
		srv: &http.Server{
			Addr:              container.Config.Server.GRPCPort,
			Handler:           h2c.NewHandler(handler, &http2.Server{}),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
	}
}

// Start starts the gRPC server
func (s *Server) Start() {
	go func() {
		s.srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// Middleware de CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Connect-Protocol-Version")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Expose-Headers", "Connect-Protocol-Version")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GrpcServerConfig holds gRPC server configuration
type GrpcServerConfig struct {
	Port              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

// DefaultGrpcServerConfig returns default gRPC server configuration
func DefaultGrpcServerConfig() GrpcServerConfig {
	return GrpcServerConfig{
		Port:              ":9000",
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

// SetupRedis sets up Redis connection
func SetupRedis(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

// SetupSMTPPort parses SMTP port from config
func SetupSMTPPort(cfg config.SMTPConfig) (int, error) {
	return strconv.Atoi(cfg.Port)
}
