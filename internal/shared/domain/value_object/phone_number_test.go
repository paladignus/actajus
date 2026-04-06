// Package valueobject
package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneNumber(t *testing.T) {
	t.Run("should valid phone number", func(t *testing.T) {
		p := PhoneNumber("+5511999999999")
		assert.True(t, p.IsValid())
		p = PhoneNumber("+556734254200")
		assert.True(t, p.IsValid())
	})
	t.Run("should invalid phone number", func(t *testing.T) {
		p := PhoneNumber("+55679934254200")
		assert.False(t, p.IsValid())
	})
}
