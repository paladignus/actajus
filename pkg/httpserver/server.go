// Package httpserver provides the HTTP server
package httpserver

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/shared/infrastructure/observability/metrics"
	"github.com/paladignus/actajus/internal/shared/infrastructure/observability/tracing"
	"github.com/paladignus/actajus/internal/shared/presentation/http/middleware"
	"github.com/paladignus/actajus/pkg/di"
)

// Server represents the HTTP server
type Server struct {
	container *di.Container
	srv       *http.Server
}

// New creates a new HTTP server
func New(container *di.Container) *Server {
	mux := http.NewServeMux()

	// Health check with dependency status
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := "ok"
		checks := map[string]any{}

		// Check PostgreSQL - need to assert to *pgxpool.Pool for Ping
		ctx := r.Context()
		if pool, ok := container.DB.(*pgxpool.Pool); ok {
			if err := pool.Ping(ctx); err != nil {
				status = "degraded"
				checks["database"] = map[string]string{"status": "unreachable", "error": err.Error()}
			} else {
				checks["database"] = map[string]string{"status": "healthy"}
			}
		}

		// Check Redis
		if container.RDB != nil {
			type pinger interface{ Ping(context.Context) interface{ Err() error } }
			if pinger, ok := container.RDB.(pinger); ok {
				if err := pinger.Ping(ctx).Err(); err != nil {
					status = "degraded"
					checks["redis"] = map[string]string{"status": "unreachable", "error": err.Error()}
				} else {
					checks["redis"] = map[string]string{"status": "healthy"}
				}
			}
		}

		jsonBytes := []byte(fmt.Sprintf(`{"status":"%s","checks":{}}`, status))
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})

	// Prometheus metrics endpoint
	mux.Handle("GET /metrics", metrics.Handler())

	// Register module routes
	container.Modules.Web.Mount(mux)

	// Build middleware chain: metrics -> tracing -> logging -> cors -> error page -> handlers
	handler := metrics.Middleware(
		tracing.Middleware("actajus-api")(
			middleware.LoggerMiddleware(container.Logger)(
				corsMiddleware(errorPageMiddleware(container, mux)),
			),
		),
	)

	return &Server{
		container: container,
		srv: &http.Server{
			Addr:         normalizeAddr(container.Config.Server.Port),
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

type statusCapturingResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	body        bytes.Buffer
}

func (w *statusCapturingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *statusCapturingResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.status = status
	w.wroteHeader = true
}

func (w *statusCapturingResponseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(p)
}

func normalizeAddr(port string) string {
	if port == "" {
		return ":8080"
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}

// Start starts the HTTP server
func (s *Server) Start() {
	// Start runtime metrics collection
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			metrics.CollectRuntimeMetrics()
		}
	}()

	// Graceful shutdown
	go func() {
		log.Printf("🚀 Server starting on %s", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
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

func errorPageMiddleware(container *di.Container, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shouldRenderHTMLErrorPage(r) {
			next.ServeHTTP(w, r)
			return
		}

		rec := &statusCapturingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(rec, r)

		status := rec.status
		if status == 0 {
			status = http.StatusOK
		}

		switch status {
		case http.StatusNotFound:
			container.Modules.Web.RenderNotFound(w, r)
			return
		case http.StatusForbidden:
			container.Modules.Web.RenderForbidden(w, r)
			return
		default:
			w.WriteHeader(status)
			_, _ = w.Write(rec.body.Bytes())
		}
	})
}

func shouldRenderHTMLErrorPage(r *http.Request) bool {
	if r.Method != http.MethodGet {
		return false
	}
	if strings.HasPrefix(r.URL.Path, "/assets/") || r.URL.Path == "/health" || strings.HasPrefix(r.URL.Path, "/bootstrap") {
		return false
	}
	accept := strings.ToLower(r.Header.Get("Accept"))
	return accept == "" || strings.Contains(accept, "text/html") || strings.Contains(accept, "*/*")
}
