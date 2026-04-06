// Package metrics provides Redis metrics for Prometheus.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Redis command metrics.
	RedisCommandCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_commands_total",
			Help: "Total number of Redis commands executed",
		},
		[]string{"command"},
	)

	RedisCommandDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "redis_command_duration_seconds",
			Help: "Duration of Redis commands in seconds",
			Buckets: []float64{
				0.0001, 0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1,
			},
		},
		[]string{"command"},
	)

	// Redis connection metrics.
	RedisConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_connections",
			Help: "Number of Redis connections (active, idle, etc.)",
		},
		[]string{"state"},
	)

	// Redis error metrics.
	RedisErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_errors_total",
			Help: "Total number of Redis errors",
		},
		[]string{"error_type"},
	)
)
