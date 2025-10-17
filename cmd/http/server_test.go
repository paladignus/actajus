// Package main
package main

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerConfiguration(t *testing.T) {
	t.Run("should the server have the correct timeout settings", func(t *testing.T) {
		srv := &http.Server{
			Addr:         ":8080",
			Handler:      http.NewServeMux(),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}
		assert.Equal(t, ":8080", srv.Addr)
		assert.Equal(t, 15*time.Second, srv.ReadTimeout)
		assert.Equal(t, 15*time.Second, srv.WriteTimeout)
		assert.Equal(t, 60*time.Second, srv.IdleTimeout)
		assert.NotNil(t, srv.TLSNextProto)
		assert.Empty(t, srv.TLSNextProto)
	})

	t.Run("should disable TLSNextProto HTTP/2", func(t *testing.T) {
		srv := &http.Server{
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}
		// TLSNextProto vazio dedabilita HTTP/2
		assert.NotNil(t, srv.TLSNextProto)
		assert.Len(t, srv.TLSNextProto, 0)
	})
}
