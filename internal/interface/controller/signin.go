// Package controller
package controller

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/usecase"
)

type SignIn struct {
	usecase usecase.SignIn
}

func NewSignIn(usecase usecase.SignIn) SignIn {
	return SignIn{usecase}
}

func (s SignIn) SignIn(ctx context.Context, req dto.AuthenticatedInput) (dto.AuthenticatedOutput, error) {
	return s.usecase.Execute(ctx, req.CPF, req.Password)
}

// 	return dto.AuthenticatedOutput{
// 		AccessToken:  "token",
// 		RefreshToken: "token",
// 		UserID:       "1",
// 		FirstName:    "John",
// 		LastName:     "Doe",
// 		Email:        "ZV5oq@example.com",
// 		Roles:        []string{"admin"},
// 	}, nil
// }
//
//
// isaac de abreu pereira, 039.809.651-12, so exclusao
