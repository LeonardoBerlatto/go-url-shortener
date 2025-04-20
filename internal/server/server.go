package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/leonardoberlatto/go-url-shortener/internal/config"
	"github.com/leonardoberlatto/go-url-shortener/internal/handlers"
	"github.com/leonardoberlatto/go-url-shortener/internal/logger"
	"github.com/leonardoberlatto/go-url-shortener/internal/middleware"
	"github.com/leonardoberlatto/go-url-shortener/internal/routes"
	"github.com/leonardoberlatto/go-url-shortener/internal/service"
	"github.com/leonardoberlatto/go-url-shortener/internal/storage"
)

type Server struct {
	router *gin.Engine
}

func Initialize(env config.Config) (*Server, error) {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	server := &Server{router: router}

	dynamoStorage, err := storage.NewDynamoDBStorage(
		env.DynamoDBEndpoint,
		env.AWSRegion,
		env.AWSAccessKeyID,
		env.AWSSecretAccessKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DynamoDB storage: %w", err)
	}

	redisCache, err := storage.NewRedisCache(env.RedisURL)
	if err != nil {
		logger.Fatalf("Error: Failed to initialize Redis cache: %v", err)
	}

	baseURL := fmt.Sprintf("%s:%s", env.Host, env.Port)
	urlService := service.NewURLService(dynamoStorage, redisCache, baseURL)

	urlHandler := handlers.NewURLHandler(urlService)

	routes.SetupRoutes(server.router, urlHandler)

	return server, nil
}

func (server *Server) Start(port string) error {
	return server.router.Run(fmt.Sprintf(":%s", port))
}
