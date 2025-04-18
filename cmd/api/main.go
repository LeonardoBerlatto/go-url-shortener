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

	server := server.Initialize()

	err = server.Start(env.Port)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
