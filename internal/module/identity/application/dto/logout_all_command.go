// Package dto
package dto

type LogoutAllCommand struct {
	IDUser int64 `json:"id_user" validate:"required|min=1"`
}
