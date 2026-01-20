// Package metrics provides Prometheus metrics for the application
// This includes HTTP request metrics, database connection metrics, and custom business metrics
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define global metrics that will be tracked across the application
var (
	// HTTP request metrics - track request count, duration, and status codes
	RequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request duration in seconds",
			Buckets: []float64{
				0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
			},
		},
		[]string{"method", "path", "status"},
	)

	// Database metrics - track database connection and query metrics
	DatabaseQueryCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "table"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "database_query_duration_seconds",
			Help: "Database query duration in seconds",
			Buckets: []float64{
				0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5,
			},
		},
		[]string{"operation", "table"},
	)

	// Business metrics - track specific business operations
	AuthSuccessCount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_success_total",
			Help: "Total number of successful authentication attempts",
		},
	)

	AuthFailureCount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_failures_total",
			Help: "Total number of failed authentication attempts",
		},
	)

	PasswordResetCount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "password_reset_requests_total",
			Help: "Total number of password reset requests",
		},
	)

	EmailSentCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "emails_sent_total",
			Help: "Total number of emails sent",
		},
		[]string{"template"},
	)
)

// Handler returns an HTTP handler for Prometheus metrics endpoint
// This handler serves the /metrics endpoint that Prometheus scrapes
func Handler() http.Handler {
	return promhttp.Handler()
}