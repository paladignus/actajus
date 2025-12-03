package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResponseWriter(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: recorder,
		statusCode:     http.StatusOK,
	}
	rw.WriteHeader(http.StatusCreated)
	assert.Equal(t, http.StatusCreated, rw.statusCode)
	data := []byte("test response")
	n, err := rw.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, len(data), rw.bytesWritten)
}

func TestLoggerMiddleware(t *testing.T) {
	logger := &spy.Logger{}
	middleware := LoggerMiddleware(logger)
	t.Run("should successful request", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		sut := middleware(nextHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		w := httptest.NewRecorder()
		logger.On("Info", req.Context(), "http request started", "method", "GET", "path", "/test", "remote_addr", "192.0.2.1:1234", "user_agent", "test-agent")
		logger.On("Info", req.Context(), "http request completed", "method", "GET", "path", "/test", "status_code", http.StatusOK, "duration_ms", int64(0), "bytes_written", mock.Anything)
		sut.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.String(), "OK")
	})

	t.Run("should return an 4xx error request", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
		})
		sut := middleware(nextHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		w := httptest.NewRecorder()
		logger.On("Warn", req.Context(), "http request completed", "method", "GET", "path", "/test", "status_code", http.StatusBadRequest, "duration_ms", int64(0), "bytes_written", mock.Anything)
		sut.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusBadRequest)
	})

	t.Run("should return an 5xx error request", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		})
		sut := middleware(nextHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		w := httptest.NewRecorder()
		logger.On("Error", req.Context(), "http request completed", "method", "GET", "path", "/test", "status_code", http.StatusInternalServerError, "duration_ms", int64(0), "bytes_written", mock.Anything)
		sut.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
	})

	t.Run("should request timing", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond) // Small delay to ensure measurable time
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Delayed response"))
		})
		sut := middleware(nextHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		w := httptest.NewRecorder()
		logger.On("Info", req.Context(), "http request completed", "method", "GET", "path", "/test", "status_code", http.StatusOK, "duration_ms", int64(10), "bytes_written", mock.Anything)
		sut.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
	})
}
