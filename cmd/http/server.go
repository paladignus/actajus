package main

import (
	"log"

	"github.com/paladignus/actajus/pkg/di"
	"github.com/paladignus/actajus/pkg/httpserver"
)

func main() {
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	server := httpserver.New(container)
	server.Start()
}
