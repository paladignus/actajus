// Package metrics provides Go runtime metrics for Prometheus.
package metrics

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Go runtime custom metrics (not provided by default).
	GoRuntimeGoroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_go_goroutines",
			Help: "Number of goroutines that currently exist",
		},
	)

	GoRuntimeMemoryAlloc = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_memory_alloc_bytes",
			Help: "Number of bytes allocated and still in use",
		},
	)

	GoRuntimeGCCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_gc_count",
			Help: "Number of garbage collections",
		},
	)
)

// CollectRuntimeMetrics collects current Go runtime metrics and updates the Prometheus metrics.
func CollectRuntimeMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	GoRuntimeGoroutines.Set(float64(runtime.NumGoroutine()))
	GoRuntimeMemoryAlloc.Set(float64(m.Alloc))
	GoRuntimeGCCount.Set(float64(m.NumGC))
}
