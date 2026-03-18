package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/infrastructure/security"
)

func main() {
	hasher := security.NewArgon2idPasswordHasher()
	hashedPassword, _ := hasher.Hash("Senha@123")
	fmt.Println(hashedPassword)

	agora := time.Now()

	// 2. Definir um tempo no futuro (ex: daqui a 1 hora e 30 minutos)
	futuro := agora.Add(1*time.Hour + 30*time.Minute)

	// 3. Calcular a diferença (Subtrair agora de futuro)
	diferenca := futuro.Sub(agora)

	// 4. Converter a diferença para minutos
	minutos := diferenca.Minutes()

	fmt.Printf("Agora: %s\n", agora.Format("15:04:05"))
	fmt.Printf("Futuro: %s\n", futuro.Format("15:04:05"))
	fmt.Printf("Minutos restantes: %.0f minutos\n", minutos)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		limit := query.Get("limit")
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		baseURL := scheme + "://" + r.Host + r.URL.Path
		fmt.Println(baseURL, limit, r.URL.Path)
	})
	http.ListenAndServe(":8081", nil)
}
