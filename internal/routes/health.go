package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardoberlatto/go-url-shortener/internal/handlers"
)

func registerHealthRoutes(router *gin.Engine) {
	router.GET("/health", handlers.Health)
}
