package main

import (
	"log"

	"github.com/leonardoberlatto/go-url-shortener/internal/config"
	"github.com/leonardoberlatto/go-url-shortener/internal/server"
)

func main() {
	env, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	server, err := server.Initialize(env)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	err = server.Start(env.Port)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

	log.Printf("Server started on port %s", env.Port)
}
