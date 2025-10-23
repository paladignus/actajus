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

func TestLoggerMiddleware_ErrorStatus(t *testing.T) {
	spyLogger := &spy.SpyLogger{}
	spyLogger.On("Info",
		mock.Anything,
		"http request started",
		"method", "GET", "path", "/error", "remote_addr", "192.0.2.1:1234", "user_agent", "",
	).Once()
	spyLogger.On("Error",
		mock.Anything,
		"http request completed",
		"method", "GET", "path", "/error", "status_code", 500, "duration_ms", mock.MatchedBy(func(duration int64) bool {
			return duration >= 0
		}), "bytes_written", 0,
	).Once()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	handler := LoggerMiddleware(spyLogger)(nextHandler)
	req := httptest.NewRequest("GET", "/error", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	spyLogger.AssertExpectations(t)
}

func TestLoggerMiddleware_WarningStatus(t *testing.T) {
	spyLogger := new(spy.SpyLogger)
	spyLogger.On("Info",
		mock.Anything,
		"http request started",
		"method", "GET", "path", "/not-found", "remote_addr", "192.0.2.1:1234", "user_agent", "",
	).Once()
	spyLogger.On("Warn",
		mock.Anything,
		"http request completed",
		"method", "GET", "path", "/not-found", "status_code", 404, "duration_ms", mock.MatchedBy(func(duration int64) bool {
			return duration >= 0
		}), "bytes_written", 0,
	).Once()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	handler := LoggerMiddleware(spyLogger)(nextHandler)
	req := httptest.NewRequest("GET", "/not-found", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	spyLogger.AssertExpectations(t)
}

func TestResponseWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: rr,
		statusCode:     http.StatusOK,
	}
	rw.WriteHeader(http.StatusCreated)
	assert.Equal(t, http.StatusCreated, rw.statusCode)
	assert.Equal(t, http.StatusCreated, rr.Code)
	data := []byte("test data")
	n, err := rw.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, len(data), rw.bytesWritten)
	assert.Equal(t, "test data", rr.Body.String())
}
