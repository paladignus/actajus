// Package metrics provides NATS messaging metrics for Prometheus.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// NATS publish metrics.
	NATSPublishCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_messages_published_total",
			Help: "Total number of messages published to NATS",
		},
		[]string{"stream", "subject"},
	)

	NATSPublishDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "nats_publish_duration_seconds",
			Help: "Duration of NATS message publish in seconds",
			Buckets: []float64{
				0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1,
			},
		},
		[]string{"stream", "subject"},
	)

	NATSPublishErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_publish_errors_total",
			Help: "Total number of NATS publish errors",
		},
		[]string{"error_type"},
	)

	// NATS consume metrics.
	NATSConsumeCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_messages_consumed_total",
			Help: "Total number of messages consumed from NATS",
		},
		[]string{"stream", "consumer"},
	)

	NATSConsumeDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "nats_consume_duration_seconds",
			Help: "Duration of NATS message consumption in seconds",
			Buckets: []float64{
				0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1,
			},
		},
		[]string{"stream", "consumer"},
	)

	NATSConsumeErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_consume_errors_total",
			Help: "Total number of NATS consume errors",
		},
		[]string{"error_type"},
	)
)
