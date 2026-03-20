// Package dto
package dto

type LogoutCommand struct {
	IDSession int64 `json:"id_session" validate:"required|min=1"`
}
