package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardoberlatto/go-url-shortener/internal/handlers"
)

func registerHealthRoutes(router *gin.Engine) {
	healthGroup := router.Group("/health")
	{
		healthGroup.GET("", handlers.Health)
	}
}
