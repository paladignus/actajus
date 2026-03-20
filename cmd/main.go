package main

import (
	"log"

	"github.com/paladignus/actajus/pkg/di"
	"github.com/paladignus/actajus/pkg/httpserver"
)

func main() {
	// Initialize DI Container
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Start HTTP server
	server := httpserver.New(container)
	server.Start()
}
