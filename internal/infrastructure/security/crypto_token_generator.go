// Package security
package security

import (
	"crypto/rand"
	"encoding/base64"
)

type CryptoTokenGenerator struct {
	length int
}

func NewCryptoTokenGenerator() CryptoTokenGenerator {
	return CryptoTokenGenerator{
		length: 32, // 32 bytes = 256 bits
	}
}

func (g CryptoTokenGenerator) Generate() (string, error) {
	bytes := make([]byte, g.length)
	_, err := rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes), err
}
