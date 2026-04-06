// Package service
package service

type MessageSerializer interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}
