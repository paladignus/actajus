// Package command
package command

type LogoutAllCommand struct {
	IDUser int64 `json:"id_user" validate:"required|min=1"`
}
