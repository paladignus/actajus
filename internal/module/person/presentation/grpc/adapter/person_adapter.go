// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/person/application/command"
	personv1 "github.com/paladignus/actajus/proto/person/v1"
)

func ProtoToCreatePersonDTO(req *personv1.CreatePersonRequest) command.CreatePersonCommand {
	return command.CreatePersonCommand{
		Name:     req.Name,
		Gender:   uint(req.Gender),
		Birthday: req.Birthday,
	}
}
