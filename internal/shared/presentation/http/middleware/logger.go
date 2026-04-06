// Package middleware provides HTTP middlewares for the application.
package middleware

import (
	"net/http"
	"time"

	"github.com/paladignus/actajus/internal/shared/application/repository"
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
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			logger.Info(r.Context(), "http request started",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
			next.ServeHTTP(rw, r)
			logLevel := logger.Info
			if rw.statusCode >= 400 && rw.statusCode < 500 {
				logLevel = logger.Warn
			} else if rw.statusCode >= 500 {
				logLevel = logger.Error
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
