package audit

import (
	"context"
	"net/http"

	sharedrepo "github.com/paladignus/actajus/internal/shared/application/repository"
)

type Adapter struct {
	logger   sharedrepo.Logger
	clientIP func(*http.Request) string
}

func New(logger sharedrepo.Logger, clientIP func(*http.Request) string) *Adapter {
	return &Adapter{
		logger:   logger,
		clientIP: clientIP,
	}
}

func (a *Adapter) Success(ctx context.Context, r *http.Request, action string, attrs ...any) {
	if a == nil || a.logger == nil || r == nil {
		return
	}
	base := []any{
		"category", "web_audit",
		"outcome", "success",
		"action", action,
		"method", r.Method,
		"path", r.URL.Path,
		"ip", a.clientIP(r),
		"user_agent", r.UserAgent(),
	}
	a.logger.Info(ctx, "web audit", append(base, attrs...)...)
}

func (a *Adapter) Failure(ctx context.Context, r *http.Request, action string, attrs ...any) {
	if a == nil || a.logger == nil || r == nil {
		return
	}
	base := []any{
		"category", "web_audit",
		"outcome", "failure",
		"action", action,
		"method", r.Method,
		"path", r.URL.Path,
		"ip", a.clientIP(r),
		"user_agent", r.UserAgent(),
	}
	a.logger.Warn(ctx, "web audit", append(base, attrs...)...)
}
