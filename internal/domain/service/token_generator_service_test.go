// Pachage service
package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ITokenGeneratorSpy struct{}

func (t ITokenGeneratorSpy) Generate() (string, error) {
	return "token", nil
}

func TestTokenGeneratorService(t *testing.T) {
	t.Run("should return success to generate token", func(t *testing.T) {
		sut := ITokenGeneratorSpy{}
		token, err := sut.Generate()
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
