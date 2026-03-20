// Package mapper
package mapper

type ChangePasswordNormalized struct {
	IDUser          int64
	CurrentPassword string
	NewPassword     string
}
