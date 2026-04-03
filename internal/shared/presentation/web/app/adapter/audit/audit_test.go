package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sharedrepo "github.com/paladignus/actajus/internal/shared/application/repository"
)

type captureLogger struct {
	level string
	msg   string
	args  []any
}

func (l *captureLogger) Debug(context.Context, string, ...any) {}
func (l *captureLogger) Info(_ context.Context, msg string, args ...any) {
	l.level = "info"
	l.msg = msg
	l.args = args
}
func (l *captureLogger) Warn(_ context.Context, msg string, args ...any) {
	l.level = "warn"
	l.msg = msg
	l.args = args
}
func (l *captureLogger) Error(context.Context, string, ...any) {}
func (l *captureLogger) With(...any) sharedrepo.Logger         { return l }
func (l *captureLogger) WithError(error) sharedrepo.Logger     { return l }

func TestAuditSuccessLogsStructuredEvent(t *testing.T) {
	logger := &captureLogger{}
	adapter := New(logger, func(*http.Request) string { return "127.0.0.1" })
	req := httptest.NewRequest(http.MethodPost, "/login", nil)

	adapter.Success(context.Background(), req, "auth.login", "id_user", int64(7))

	if logger.level != "info" || logger.msg != "web audit" {
		t.Fatalf("unexpected log entry: %s %s", logger.level, logger.msg)
	}
	if len(logger.args) == 0 {
		t.Fatalf("expected structured args")
	}
}

func TestAuditFailureLogsStructuredEvent(t *testing.T) {
	logger := &captureLogger{}
	adapter := New(logger, func(*http.Request) string { return "127.0.0.1" })
	req := httptest.NewRequest(http.MethodPost, "/login", nil)

	adapter.Failure(context.Background(), req, "auth.login", "reason", "invalid_credentials")

	if logger.level != "warn" || logger.msg != "web audit" {
		t.Fatalf("unexpected log entry: %s %s", logger.level, logger.msg)
	}
}
