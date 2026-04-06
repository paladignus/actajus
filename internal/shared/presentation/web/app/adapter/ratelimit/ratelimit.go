package ratelimit

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type entry struct {
	failures     int
	firstFailure time.Time
	lockedUntil  time.Time
}

type LoginShield struct {
	mu          sync.Mutex
	now         func() time.Time
	window      time.Duration
	threshold   int
	baseBackoff time.Duration
	maxBackoff  time.Duration
	clientIP    func(*http.Request) string
	byIP        map[string]entry
	byEmail     map[string]entry
}

func NewLoginShield(clientIP func(*http.Request) string) *LoginShield {
	if clientIP == nil {
		clientIP = func(*http.Request) string { return "" }
	}
	return &LoginShield{
		now:         time.Now,
		window:      15 * time.Minute,
		threshold:   3,
		baseBackoff: 30 * time.Second,
		maxBackoff:  15 * time.Minute,
		clientIP:    clientIP,
		byIP:        map[string]entry{},
		byEmail:     map[string]entry{},
	}
}

func (s *LoginShield) Check(r *http.Request, email string) (bool, time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	ipRetry := s.retryAfter(now, s.byIP, s.clientIP(r))
	emailRetry := s.retryAfter(now, s.byEmail, normalizeEmail(email))
	retryAfter := maxDuration(ipRetry, emailRetry)
	return retryAfter == 0, retryAfter
}

func (s *LoginShield) RecordFailure(r *http.Request, email string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	s.recordFailure(now, s.byIP, s.clientIP(r))
	s.recordFailure(now, s.byEmail, normalizeEmail(email))
}

func (s *LoginShield) RecordSuccess(r *http.Request, email string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.byIP, s.clientIP(r))
	delete(s.byEmail, normalizeEmail(email))
}

func (s *LoginShield) retryAfter(now time.Time, store map[string]entry, key string) time.Duration {
	if key == "" {
		return 0
	}
	item, ok := store[key]
	if !ok {
		return 0
	}
	item = s.prune(now, item)
	if item.failures == 0 {
		delete(store, key)
		return 0
	}
	store[key] = item
	if item.lockedUntil.After(now) {
		return item.lockedUntil.Sub(now)
	}
	return 0
}

func (s *LoginShield) recordFailure(now time.Time, store map[string]entry, key string) {
	if key == "" {
		return
	}
	item := s.prune(now, store[key])
	if item.failures == 0 {
		item.firstFailure = now
	}
	item.failures++

	if item.failures >= s.threshold {
		lockDuration := s.baseBackoff << (item.failures - s.threshold)
		if lockDuration > s.maxBackoff {
			lockDuration = s.maxBackoff
		}
		item.lockedUntil = now.Add(lockDuration)
	}
	store[key] = item
}

func (s *LoginShield) prune(now time.Time, item entry) entry {
	if item.failures == 0 {
		return entry{}
	}
	if item.firstFailure.IsZero() || now.Sub(item.firstFailure) > s.window {
		return entry{}
	}
	if item.lockedUntil.Before(now) {
		item.lockedUntil = time.Time{}
	}
	return item
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
