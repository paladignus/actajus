// Package mapper
package mapper

type ConfirmPasswordResetNormalized struct {
	IDReset     int64
	ResetToken  string
	NewPassword string
}
