package main

import (
	"github.com/leonardoberlatto/go-url-shortener/internal/config"
	"github.com/leonardoberlatto/go-url-shortener/internal/logger"
	"github.com/leonardoberlatto/go-url-shortener/internal/server"
)

func main() {
	defer logger.Sync()
	logger.Init(logger.InfoLevel)

	env, err := config.Load()
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	logger.Infof("Starting application with log level: %s", env.LogLevel)

	server, err := server.Initialize(env)
	if err != nil {
		logger.Fatalf("Error initializing server: %v", err)
	}

	err = server.Start(env.Port)
	if err != nil {
		logger.Fatal("Cannot start server:", err)
	}
}
