// Package valueobject
package valueobject

type PasswordHash string

func (p PasswordHash) Value() string { return string(p) }
func (p PasswordHash) IsEmpty() bool { return len(p) == 0 }
