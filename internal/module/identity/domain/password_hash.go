// Package domain
package domain

type PasswordHash string

func (p PasswordHash) Value() string { return string(p) }
func (p PasswordHash) IsEmpty() bool { return len(p) == 0 }
