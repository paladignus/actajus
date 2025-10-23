// Package middleware
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnableCORS(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   bool
	}{
		{
			name:           "OPTIONS request - preflight",
			method:         "OPTIONS",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "GET request",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "POST request",
			method:         "POST",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "PUT request",
			method:         "PUT",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			handler := EnableCORS(nextHandler)
			req := httptest.NewRequest(tt.method, "/", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.checkHeaders {
				assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
			}
		})
	}
}

func TestEnableCORS_OptionsRequestStopsPropagation(t *testing.T) {
	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handler := EnableCORS(nextHandler)
	req := httptest.NewRequest("OPTIONS", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.False(t, called, "Next handler should not be called for OPTIONS requests")
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestEnableCORS_NonOptionsRequestContinues(t *testing.T) {
	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handler := EnableCORS(nextHandler)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.True(t, called, "Next handler should be called for non-OPTIONS requests")
	assert.Equal(t, http.StatusOK, rr.Code)
}
