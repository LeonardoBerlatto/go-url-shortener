package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardoberlatto/go-url-shortener/internal/handlers"
)

func SetupRoutes(router *gin.Engine, urlHandler *handlers.URLHandler) {
	registerHealthRoutes(router)
	registerURLRoutes(router, urlHandler)
}
