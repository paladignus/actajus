// Package metrics provides Redis metrics hook for Prometheus.
package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

// RedisMetricsHook implements redis.Hook interface to collect metrics.
type RedisMetricsHook struct{}

// NewRedisMetricsHook creates a new Redis metrics hook.
func NewRedisMetricsHook() *RedisMetricsHook {
	return &RedisMetricsHook{}
}

// DialHook implements redis.Hook.
func (h *RedisMetricsHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

// ProcessHook implements redis.Hook - collects metrics for each command.
func (h *RedisMetricsHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		duration := time.Since(start)

		commandName := cmd.Name()
		
		// Record command count
		RedisCommandCount.With(prometheus.Labels{"command": commandName}).Inc()
		
		// Record command duration
		RedisCommandDuration.With(prometheus.Labels{"command": commandName}).Observe(duration.Seconds())
		
		// Record errors
		if err != nil {
			RedisErrors.With(prometheus.Labels{"error_type": commandName}).Inc()
		}
		
		return err
	}
}

// ProcessPipelineHook implements redis.Hook.
func (h *RedisMetricsHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}
