// Package valueobject
package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	urls := []string{
		"https://github.com/paladignus/actajus",
		"https://x.com/actajus",
		"https://facebook.com/actajus",
		"https://instagram.com/actajus",
	}
	t.Run("should true if url are valid", func(t *testing.T) {
		for _, url := range urls {
			assert.Equal(t, true, URL(url).IsValid(), url)
			assert.True(t, URL(url).IsValid(), true)
		}
	})
	t.Run("should false if url are invalid", func(t *testing.T) {
		invalidUrls := []string{
			"invalid-url",
			"invalid-url.com",
			"invalid-url.com.br",
			"invalid-url.com.br/actajus",
		}
		for _, url := range invalidUrls {
			assert.Equal(t, false, URL(url).IsValid(), url)
			assert.False(t, URL(url).IsValid(), false)
		}
	})
}
