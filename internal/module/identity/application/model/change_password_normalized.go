// Package model
package model

type ChangePasswordNormalized struct {
	IDUser          int64
	CurrentPassword string
	NewPassword     string
}
