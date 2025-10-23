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
