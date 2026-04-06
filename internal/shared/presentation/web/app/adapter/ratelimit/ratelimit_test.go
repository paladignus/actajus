package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoginShieldAllowsBeforeThreshold(t *testing.T) {
	now := time.Unix(100, 0)
	shield := NewLoginShield(func(*http.Request) string { return "127.0.0.1" })
	shield.now = func() time.Time { return now }

	req := httptest.NewRequest(http.MethodPost, "/login", nil)

	for range 2 {
		shield.RecordFailure(req, "user@mail.com")
	}

	allowed, retryAfter := shield.Check(req, "user@mail.com")
	if !allowed || retryAfter != 0 {
		t.Fatalf("expected request allowed before threshold, got allowed=%v retry=%v", allowed, retryAfter)
	}
}

func TestLoginShieldBlocksAfterThresholdByEmailAndIP(t *testing.T) {
	now := time.Unix(100, 0)
	shield := NewLoginShield(func(*http.Request) string { return "127.0.0.1" })
	shield.now = func() time.Time { return now }

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	for range 3 {
		shield.RecordFailure(req, "USER@mail.com")
	}

	allowed, retryAfter := shield.Check(req, "user@mail.com")
	if allowed {
		t.Fatalf("expected request blocked")
	}
	if retryAfter <= 0 {
		t.Fatalf("expected retry after > 0, got %v", retryAfter)
	}
}

func TestLoginShieldBackoffGrowsWithRepeatedFailures(t *testing.T) {
	now := time.Unix(100, 0)
	shield := NewLoginShield(func(*http.Request) string { return "127.0.0.1" })
	shield.now = func() time.Time { return now }

	req := httptest.NewRequest(http.MethodPost, "/login", nil)

	for range 3 {
		shield.RecordFailure(req, "user@mail.com")
	}
	_, firstRetry := shield.Check(req, "user@mail.com")

	now = now.Add(firstRetry + time.Second)
	shield.RecordFailure(req, "user@mail.com")
	_, secondRetry := shield.Check(req, "user@mail.com")

	if secondRetry <= firstRetry {
		t.Fatalf("expected backoff to grow, first=%v second=%v", firstRetry, secondRetry)
	}
}

func TestLoginShieldRecordSuccessClearsLock(t *testing.T) {
	now := time.Unix(100, 0)
	shield := NewLoginShield(func(*http.Request) string { return "127.0.0.1" })
	shield.now = func() time.Time { return now }

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	for range 3 {
		shield.RecordFailure(req, "user@mail.com")
	}

	shield.RecordSuccess(req, "user@mail.com")
	allowed, retryAfter := shield.Check(req, "user@mail.com")
	if !allowed || retryAfter != 0 {
		t.Fatalf("expected success to clear shield, got allowed=%v retry=%v", allowed, retryAfter)
	}
}
