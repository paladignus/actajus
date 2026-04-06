// Package handler
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type InputUserEmail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestDecodeJSONRequest(t *testing.T) {
	t.Run("should valid JSON request", func(t *testing.T) {
		type NameEmailInput struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		input := NameEmailInput{
			Name:  "John Doe",
			Email: "john@email.com",
		}
		bodyBytes, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("failed to marshal input: %v", err)
		}
		req := httptest.NewRequest("POST", "/test", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		result, err := DecodeJSONRequest[NameEmailInput](req)
		assert.NoError(t, err)
		expected, _ := json.Marshal(result)
		assert.JSONEq(t, string(expected), string(bodyBytes))
	})

	t.Run("should invalid JSON request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		_, err := DecodeJSONRequest[InputUserEmail](req)
		assert.Error(t, err)
	})

	t.Run("should empty JSON request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		_, err := DecodeJSONRequest[InputUserEmail](req)
		assert.Error(t, err)
	})
}

func TestRespondJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       any
	}{
		{
			name:       "should perform the coding of a simple structure",
			statusCode: http.StatusOK,
			data: struct {
				Message string `json:"message"`
			}{
				Message: "Hello World",
			},
		},
		{
			name:       "should perform the coding of a complex structure",
			statusCode: http.StatusCreated,
			data: struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			}{
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
		{
			name:       "should perform the coding of a slice data",
			statusCode: http.StatusOK,
			data:       []string{"apple", "banana", "orange"},
		},
		{
			name:       "should not be coded for empty data",
			statusCode: http.StatusNoContent,
			data:       nil,
		},
		{
			name:       "should do the coding for a map",
			statusCode: http.StatusOK,
			data: map[string]any{
				"success": true,
				"count":   42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			err := RespondJSON(resp, tt.statusCode, tt.data)
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.Code)
			assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
			if tt.data != nil {
				var responseData any
				err := json.Unmarshal(resp.Body.Bytes(), &responseData)
				assert.NoError(t, err)
				expectedBytes, _ := json.Marshal(tt.data)
				actualBytes, _ := json.Marshal(responseData)
				assert.JSONEq(t, string(expectedBytes), string(actualBytes))
			}
		})
	}
}

func TestRespondJSONError(t *testing.T) {
	t.Run("should return an error when performing the encoding", func(t *testing.T) {
		resp := httptest.NewRecorder()
		type Circular struct {
			Self *Circular `json:"self"`
		}
		circular := &Circular{}
		circular.Self = circular
		err := RespondJSON(resp, http.StatusOK, circular)
		assert.Error(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
	})
}

func TestRespondError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "should return bad request",
			statusCode: http.StatusBadRequest,
			message:    "Invalid input",
		},
		{
			name:       "should return not found",
			statusCode: http.StatusNotFound,
			message:    "Resource not found",
		},
		{
			name:       "should return internal server error",
			statusCode: http.StatusInternalServerError,
			message:    "Something went wrong",
		},
		{
			name:       "should return unauthorized",
			statusCode: http.StatusUnauthorized,
			message:    "Access denied",
		},
		{
			name:       "should return empty message",
			statusCode: http.StatusBadRequest,
			message:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			RespondError(resp, tt.statusCode, tt.message)
			assert.Equal(t, tt.statusCode, resp.Code)
			assert.Equal(t, tt.message+"\n", resp.Body.String())
			assert.Equal(t, "text/plain; charset=utf-8", resp.Header().Get("Content-Type"))
		})
	}
}
