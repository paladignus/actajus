// Package dto
package dto

type RevokeSessionCommand struct {
	IDSession int64  `json:"id_session" validate:"required|min=1"`
	Reason    string `json:"reason"`
}
