// Package command
package command

import "time"

type CreateUserCommand struct {
	Password    string
	Avatar      string
	LastLoginAt time.Time
}
