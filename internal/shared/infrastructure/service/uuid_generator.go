// Package service
package service

import "github.com/google/uuid"

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) NewString() string {
	return uuid.NewString()
}
