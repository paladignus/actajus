package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paladignus/actajus/internal/module/company"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/internal/shared/infrastructure/logger"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

func main() {
	// ctx := context.Background()

	// // Database
	// dbConfig := database.PostgresConfig{
	// 	Host:     getEnv("DB_HOST", "localhost"),
	// 	Port:     5432,
	// 	User:     getEnv("DB_USER", "postgres"),
	// 	Password: getEnv("DB_PASSWORD", "M4rc3l0"),
	// 	Database: getEnv("DB_NAME", "actajus"),
	// 	SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	// }

	// pool, err := database.NewPostgresPool(ctx, dbConfig)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer pool.Close()
	config := config.Load()
	ctx := context.Background()
	logger := logger.NewDefaultLogger()
	db, err := postgres.NewConnection(ctx, &config.Database)
	if err != nil {
		logger.Error(ctx, "error initializing the database connection.", "error", err)
		os.Exit(1)
	}
	// defer db.Close(ctx, logger)
	defer db.Close()
	// persistence := postgres.NewPersistence(db)
	// persistence := postgres.NewPersistence(db)

	log.Println("✅ Database connected successfully")

	// Initialize modules
	// companyMod := companyModule.NewModule(pool)
	// addressMod := addressModule.NewModule(pool)

	companyModule := company.NewModule(db)

	log.Println("✅ Modules initialized")

	// ServeMux
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Register module routes
	companyModule.Handler.RegisterRoutes(mux)
	// addressMod.Handler.RegisterRoutes(mux)

	// Middleware wrapper
	handler := loggingMiddleware(corsMiddleware(mux))

	// Server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Middleware de logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call next handler
		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// Middleware de CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
