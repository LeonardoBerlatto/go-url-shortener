package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardoberlatto/go-url-shortener/internal/handlers"
)

func registerURLRoutes(router *gin.Engine, urlHandler *handlers.URLHandler) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/urls", urlHandler.ShortenURL)
		v1.DELETE("/urls/:code", urlHandler.DeleteURL)
		v1.GET("/urls", urlHandler.ListURLs)
	}

	router.GET("/:code", urlHandler.RedirectURL)
}
