package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leonardoberlatto/go-url-shortener/internal/models"
	"github.com/leonardoberlatto/go-url-shortener/internal/service"
	"github.com/leonardoberlatto/go-url-shortener/internal/storage"
)

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}

	response, err := h.urlService.Shorten(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == storage.ErrorConflict {
			statusCode = http.StatusConflict
			c.JSON(statusCode, models.ErrorResponse{Error: "Custom short ID already exists"})
			return
		}
		c.JSON(statusCode, models.ErrorResponse{Error: "Failed to shorten URL: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortID := c.Param("code")

	longURL, err := h.urlService.Resolve(c.Request.Context(), shortID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == storage.ErrorNotFound {
			statusCode = http.StatusNotFound
			c.JSON(statusCode, models.ErrorResponse{Error: "Short URL not found"})
			return
		}
		c.JSON(statusCode, models.ErrorResponse{Error: "Failed to resolve URL: " + err.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
}

func (h *URLHandler) DeleteURL(c *gin.Context) {
	shortID := c.Param("code")

	err := h.urlService.Delete(c.Request.Context(), shortID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == storage.ErrorNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, models.ErrorResponse{Error: "Failed to delete URL: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
