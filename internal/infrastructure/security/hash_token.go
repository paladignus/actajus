// Package security
package security

import (
	"crypto/rand"
	"encoding/base64"
)

type HashToken struct {
	length int
}

func NewHashToken() *HashToken {
	return &HashToken{
		length: 32, // 32 bytes = 256 bits
	}
}

func (g *HashToken) GenerateToken() (string, error) {
	bytes := make([]byte, g.length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
