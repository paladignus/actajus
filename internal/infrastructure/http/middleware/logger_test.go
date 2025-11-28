package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/paladignus/actajus/internal/domain/repository"
)

// TestResponseWriter tests the responseWriter wrapper
func TestResponseWriter(t *testing.T) {
	// Create a mock response writer
	recorder := httptest.NewRecorder()
	
	// Wrap it with our responseWriter
	rw := &responseWriter{
		ResponseWriter: recorder,
		statusCode:     http.StatusOK,
	}
	
	// Test WriteHeader
	rw.WriteHeader(http.StatusCreated)
	
	if rw.statusCode != http.StatusCreated {
		t.Errorf("Expected statusCode %d, got %d", http.StatusCreated, rw.statusCode)
	}
	
	// Test Write
	data := []byte("test response")
	n, err := rw.Write(data)
	
	if err != nil {
		t.Errorf("Write returned error: %v", err)
	}
	
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}
	
	if rw.bytesWritten != len(data) {
		t.Errorf("Expected bytesWritten %d, got %d", len(data), rw.bytesWritten)
	}
}

// MockLogger is a mock implementation of repository.Logger for testing
type MockLogger struct {
	logCalls []string
	logArgs  []interface{}
}

func (m *MockLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	m.logCalls = append(m.logCalls, "debug")
	m.logArgs = append(m.logArgs, args)
}

func (m *MockLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	m.logCalls = append(m.logCalls, "info")
	m.logArgs = append(m.logArgs, args)
}

func (m *MockLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	m.logCalls = append(m.logCalls, "warn")
	m.logArgs = append(m.logArgs, args)
}

func (m *MockLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	m.logCalls = append(m.logCalls, "error")
	m.logArgs = append(m.logArgs, args)
}

func (m *MockLogger) With(args ...interface{}) repository.Logger {
	return m
}

func (m *MockLogger) WithError(err error) repository.Logger {
	return m
}

// TestLoggerMiddleware tests the LoggerMiddleware functionality
func TestLoggerMiddleware(t *testing.T) {
	
	t.Run("successful request", func(t *testing.T) {
		mockLogger := &MockLogger{
			logCalls: []string{},
			logArgs:  []interface{}{},
		}
		
		middleware := LoggerMiddleware(mockLogger)
		
		// Create a simple handler that returns 200
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		
		handler := middleware(nextHandler)
		
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)
		
		// Should have 2 log calls: one for request start, one for completion
		if len(mockLogger.logCalls) < 2 {
			t.Errorf("Expected at least 2 log calls, got %d", len(mockLogger.logCalls))
		}
		
		// Check response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		
		if w.Body.String() != "OK" {
			t.Errorf("Expected body 'OK', got '%s'", w.Body.String())
		}
	})
	
	t.Run("4xx error request", func(t *testing.T) {
		mockLogger := &MockLogger{
			logCalls: []string{},
			logArgs:  []interface{}{},
		}
		
		middleware := LoggerMiddleware(mockLogger)
		
		// Create a handler that returns 400
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
		})
		
		handler := middleware(nextHandler)
		
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)
		
		// Should have at least 2 log calls
		if len(mockLogger.logCalls) < 2 {
			t.Errorf("Expected at least 2 log calls, got %d", len(mockLogger.logCalls))
		}
		
		// Check response
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
	
	t.Run("5xx error request", func(t *testing.T) {
		mockLogger := &MockLogger{
			logCalls: []string{},
			logArgs:  []interface{}{},
		}
		
		middleware := LoggerMiddleware(mockLogger)
		
		// Create a handler that returns 500
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		})
		
		handler := middleware(nextHandler)
		
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)
		
		// Should have at least 2 log calls
		if len(mockLogger.logCalls) < 2 {
			t.Errorf("Expected at least 2 log calls, got %d", len(mockLogger.logCalls))
		}
		
		// Check response
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
	
	t.Run("request timing", func(t *testing.T) {
		mockLogger := &MockLogger{
			logCalls: []string{},
			logArgs:  []interface{}{},
		}
		
		middleware := LoggerMiddleware(mockLogger)
		
		// Create a handler that takes some time
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond) // Small delay to ensure measurable time
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Delayed response"))
		})
		
		handler := middleware(nextHandler)
		
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)
		
		// Check that the response was processed
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}