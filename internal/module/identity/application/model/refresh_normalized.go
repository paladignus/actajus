// Package model
package model

type RefreshNormalized struct {
	IDSession    int64
	RefreshToken string
	IP           string
	UserAgent    string
}
