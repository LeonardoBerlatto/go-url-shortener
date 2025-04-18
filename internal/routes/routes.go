package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	registerHealthRoutes(router)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/shorten", nil)
		v1.GET("/:code", nil)
		v1.GET("/urls", nil)
		v1.DELETE("/:code", nil)
	}
}
