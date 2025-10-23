// Package middleware
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paladignus/actajus/test/infrastructure/adapter/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoggerMiddleware_Success(t *testing.T) {
	spyLogger := new(spy.SpyLogger)
	spyLogger.On("Info",
		mock.Anything,
		"http request started",
		"method", "GET", "path", "/test", "remote_addr", "192.168.1.1:12345", "user_agent", "test-agent",
	).Once()
	spyLogger.On("Info",
		mock.Anything,
		"http request completed",
		"method", "GET", "path", "/test", "status_code", 200, "duration_ms", mock.MatchedBy(func(duration int64) bool {
			return duration >= 0
		}), "bytes_written", 11,
	).Once()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})
	handler := LoggerMiddleware(spyLogger)(nextHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("User-Agent", "test-agent")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello World", rr.Body.String())
	spyLogger.AssertExpectations(t)
}
