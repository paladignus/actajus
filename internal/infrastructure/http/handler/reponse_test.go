// Package handler
package handler

import (
	"bytes"
	"encoding/json"
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
