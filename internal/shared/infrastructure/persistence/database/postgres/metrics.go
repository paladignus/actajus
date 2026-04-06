// Package postgres provides PostgreSQL infrastructure with observability.
package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

// PoolMetricsCollector collects Prometheus metrics for a PostgreSQL connection pool.
type PoolMetricsCollector struct {
	pool *pgxpool.Pool
	name string

	// Metrics descriptors
	totalConns          *prometheus.Desc
	acquiredConns       *prometheus.Desc
	idleConns           *prometheus.Desc
	maxConns            *prometheus.Desc
	constructingConns   *prometheus.Desc
	acquireCount        *prometheus.Desc
	emptyAcquireCount   *prometheus.Desc
	canceledAcquireCount *prometheus.Desc
}

// NewPoolMetricsCollector creates a new pool metrics collector.
func NewPoolMetricsCollector(pool *pgxpool.Pool, name string) *PoolMetricsCollector {
	labels := prometheus.Labels{"pool": name}
	return &PoolMetricsCollector{
		pool: pool,
		name: name,
		totalConns: prometheus.NewDesc(
			"database_pool_total_conns",
			"Total number of connections in the pool",
			nil, labels,
		),
		acquiredConns: prometheus.NewDesc(
			"database_pool_acquired_conns",
			"Number of currently acquired (in-use) connections",
			nil, labels,
		),
		idleConns: prometheus.NewDesc(
			"database_pool_idle_conns",
			"Number of idle connections in the pool",
			nil, labels,
		),
		maxConns: prometheus.NewDesc(
			"database_pool_max_conns",
			"Maximum number of connections allowed in the pool",
			nil, labels,
		),
		constructingConns: prometheus.NewDesc(
			"database_pool_constructing_conns",
			"Number of connections currently being established",
			nil, labels,
		),
		acquireCount: prometheus.NewDesc(
			"database_pool_acquire_count",
			"Total number of connections acquired from the pool",
			nil, labels,
		),
		emptyAcquireCount: prometheus.NewDesc(
			"database_pool_empty_acquire_count",
			"Total number of empty acquires (had to wait for connection)",
			nil, labels,
		),
		canceledAcquireCount: prometheus.NewDesc(
			"database_pool_canceled_acquire_count",
			"Total number of canceled acquires",
			nil, labels,
		),
	}
}

// Describe implements prometheus.Collector.
func (c *PoolMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalConns
	ch <- c.acquiredConns
	ch <- c.idleConns
	ch <- c.maxConns
	ch <- c.constructingConns
	ch <- c.acquireCount
	ch <- c.emptyAcquireCount
	ch <- c.canceledAcquireCount
}

// Collect implements prometheus.Collector.
func (c *PoolMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	if c.pool == nil {
		return
	}

	stats := c.pool.Stat()

	ch <- prometheus.MustNewConstMetric(
		c.totalConns, prometheus.GaugeValue, float64(stats.TotalConns()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.acquiredConns, prometheus.GaugeValue, float64(stats.AcquiredConns()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.idleConns, prometheus.GaugeValue, float64(stats.IdleConns()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.maxConns, prometheus.GaugeValue, float64(stats.MaxConns()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.constructingConns, prometheus.GaugeValue, float64(stats.ConstructingConns()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.acquireCount, prometheus.CounterValue, float64(stats.AcquireCount()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.emptyAcquireCount, prometheus.CounterValue, float64(stats.EmptyAcquireCount()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.canceledAcquireCount, prometheus.CounterValue, float64(stats.CanceledAcquireCount()),
	)
}
