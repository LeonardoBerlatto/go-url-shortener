package main

import (
	_ "github.com/leonardoberlatto/go-url-shortener/docs" // Import swagger docs
	"github.com/leonardoberlatto/go-url-shortener/internal/config"
	"github.com/leonardoberlatto/go-url-shortener/internal/logger"
	"github.com/leonardoberlatto/go-url-shortener/internal/server"
)

// @title           URL Shortener API
// @version         1.0
// @description     A simple URL shortener service API
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

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
