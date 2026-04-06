// Package service
package service

type IDGenerator interface {
	NewString() string
}
