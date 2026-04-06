// Package command
package command

type RegisterCommand struct {
	FirstName string `json:"first_name" validate:"required|min=2"`
	LastName  string `json:"last_name" validate:"required|min=2"`
	Birthday  string `json:"birthday" validate:"required"`
	GenderID  uint   `json:"gender_id" validate:"required|min=1"`
	Email     string `json:"email" validate:"required|email"`
	Password  string `json:"password" validate:"required|min=8|password"`
}
