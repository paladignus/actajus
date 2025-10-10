// Package middleware
package middleware

import (
	"net/http"
	"time"

	"github.com/paladignus/actajus/internal/domain/repository"
	appctx "github.com/paladignus/actajus/internal/infrastructure/context"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

func LoggerMiddleware(logger repository.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			traceID := appctx.GetTraceID(r.Context())
			contextLogger := logger.With("trace_id", traceID)
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			// Log da requisição
			logger.Info(r.Context(), "http request started",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
			next.ServeHTTP(rw, r)
			// Log da resposta
			logLevel := contextLogger.Info
			if rw.statusCode >= 400 && rw.statusCode < 500 {
				logLevel = contextLogger.Warn
			} else if rw.statusCode >= 500 {
				logLevel = contextLogger.Error
			}
			logLevel(r.Context(), "http request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status_code", rw.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
				"bytes_written", rw.bytesWritten,
			)
		})
	}
}
