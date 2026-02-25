package main

import (
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/infrastructure/security"
)

func main() {
	hasher := security.NewArgon2idPasswordHasher()
	hashedPassword, _ := hasher.Hash("Senha@123")
	fmt.Println(hashedPassword)
}
