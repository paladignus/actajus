// Package handler
package handler

import (
	"fmt"
	"net/http"
)

type RenewPassword struct{}

// 	RenewPasswordService service.RenewPassword
// 	Logger               repository.Logger
// }

func NewRenewPassword() RenewPassword {
	return RenewPassword{}
	// return RenewPassword{RenewPasswordService: renewPasswordService, Logger: logger}
}

func (rp RenewPassword) RenewPassword(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	fmt.Println("TOKEN=", token)
}
