// Package model
package model

type ConfirmPasswordResetNormalized struct {
	IDReset     int64
	ResetToken  string
	NewPassword string
}
