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

// ShortenURL godoc
// @Summary      Shorten a URL
// @Description  Create a short URL from a long URL
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        request  body      models.ShortenRequest  true  "URL to shorten"
// @Success      201      {object}  models.ShortenResponse
// @Failure      400      {object}  models.ErrorResponse
// @Failure      409      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /api/v1/urls [post]
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

// RedirectURL godoc
// @Summary      Redirect to original URL
// @Description  Redirects to the original URL associated with the short code
// @Tags         urls
// @Produce      plain
// @Param        code    path      string  true  "Short URL code"
// @Success      301     {string}  string  "Redirect to original URL"
// @Failure      404     {object}  models.ErrorResponse
// @Failure      500     {object}  models.ErrorResponse
// @Router       /{code} [get]
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

// DeleteURL godoc
// @Summary      Delete a shortened URL
// @Description  Deletes a shortened URL by its code
// @Tags         urls
// @Produce      json
// @Param        code    path      string  true  "Short URL code"
// @Success      204     {string}  string  "No content"
// @Failure      404     {object}  models.ErrorResponse
// @Failure      500     {object}  models.ErrorResponse
// @Router       /api/v1/urls/{code} [delete]
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

// ListURLs godoc
// @Summary      List all shortened URLs
// @Description  Returns a paginated list of all shortened URLs
// @Tags         urls
// @Produce      json
// @Param        pageNumber    query     int  false  "Page number"  default(1)  minimum(1)
// @Param        pageSize      query     int  false  "Page size"    default(10) minimum(1) maximum(100)
// @Success      200           {object}  models.PaginatedURLsResponse
// @Failure      400           {object}  models.ErrorResponse
// @Failure      500           {object}  models.ErrorResponse
// @Router       /api/v1/urls [get]
func (h *URLHandler) ListURLs(c *gin.Context) {
	var pagination models.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters: " + err.Error()})
		return
	}

	response, err := h.urlService.ListURLs(c.Request.Context(), pagination.PageNumber, pagination.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list URLs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
