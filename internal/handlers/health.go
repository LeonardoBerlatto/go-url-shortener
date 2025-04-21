package handlers

import "github.com/gin-gonic/gin"

type HealthResponse struct {
	Status string `json:"status"`
}

// Health godoc
// @Summary      Health check endpoint
// @Description  Returns the health status of the application
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func Health(context *gin.Context) {
	context.JSON(200, HealthResponse{
		Status: "ok",
	})
}
