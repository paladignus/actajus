// Package repository
package repository

import "time"

type Cache interface {
	Set(string, string, time.Duration) error
	Get(string) (string, error)
	Exists(string) (int64, error)
	Delete(string) error
}
