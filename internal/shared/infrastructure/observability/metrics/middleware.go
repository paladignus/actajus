// Package metrics provides HTTP middleware for collecting Prometheus metrics.
package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// responseWriter wrapper to capture response status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Middleware creates an HTTP middleware that collects Prometheus metrics.
// This middleware tracks:
// 1. Request count by method, path, and status code.
// 2. Request duration by method, path, and status code.
// 3. It should be placed early in the middleware chain to capture all requests.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture the status code.
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200
		}

		// Call the next handler.
		next.ServeHTTP(wrapped, r)

		// Calculate request duration.
		duration := time.Since(start)

		// Increment request counter.
		RequestCount.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": strconv.Itoa(wrapped.statusCode),
		}).Inc()

		// Observe request duration.
		RequestDuration.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": strconv.Itoa(wrapped.statusCode),
		}).Observe(duration.Seconds())
	})
}
