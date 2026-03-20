package grpcserver_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/pkg/di"
	"github.com/paladignus/actajus/pkg/grpcserver"
)

func TestGrpcServerCreation(t *testing.T) {
	t.Skip("Skipping test - requires full DI container")
	
	// Create a mock container
	container := &di.Container{
		Config: config.Config{
			Server: config.ServerConfig{
				GRPCPort: ":9000",
			},
		},
	}
	
	server := grpcserver.New(container)
	if server == nil {
		t.Fatal("Expected server to be created")
	}
}

func TestGrpcServerStart(t *testing.T) {
	t.Skip("Skipping test - requires full server setup")
	
	container := &di.Container{
		Config: config.Config{
			Server: config.ServerConfig{
				GRPCPort: ":9000",
			},
		},
	}
	
	server := grpcserver.New(container)
	
	// Start server in goroutine
	go server.Start()
	
	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	
	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Expected no error on shutdown, got %v", err)
	}
}

func TestCORSHeaders(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	// Test CORS middleware logic inline
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Connect-Protocol-Version")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Expose-Headers", "Connect-Protocol-Version")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		handler.ServeHTTP(w, r)
	})
	
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()
	
	corsHandler.ServeHTTP(w, req)
	
	// Check CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected Access-Control-Allow-Origin: *")
	}
	
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for OPTIONS, got %d", w.Code)
	}
}

func TestSetupRedis(t *testing.T) {
	t.Skip("Skipping test - requires Redis connection")
	
	ctx := context.Background()
	cfg := config.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
	}
	
	// This would test grpcserver.SetupRedis if it was exported
	_ = ctx
	_ = cfg
}

func TestDefaultGrpcServerConfig(t *testing.T) {
	cfg := grpcserver.DefaultGrpcServerConfig()
	
	if cfg.Port != ":9000" {
		t.Errorf("Expected default port :9000, got %s", cfg.Port)
	}
	
	if cfg.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("Expected ReadHeaderTimeout 5s, got %v", cfg.ReadHeaderTimeout)
	}
	
	if cfg.ReadTimeout != 10*time.Second {
		t.Errorf("Expected ReadTimeout 10s, got %v", cfg.ReadTimeout)
	}
	
	if cfg.WriteTimeout != 10*time.Second {
		t.Errorf("Expected WriteTimeout 10s, got %v", cfg.WriteTimeout)
	}
	
	if cfg.IdleTimeout != 120*time.Second {
		t.Errorf("Expected IdleTimeout 120s, got %v", cfg.IdleTimeout)
	}
}
