package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leonardoberlatto/go-url-shortener/internal/logger"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// After request
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		if statusCode >= 400 {
			logger.With(
				"status", statusCode,
				"method", method,
				"path", path,
				"latency", latency,
				"error", errorMessage,
			).Errorf("HTTP %s %s", method, path)
		} else {
			logger.With(
				"status", statusCode,
				"method", method,
				"path", path,
				"latency", latency,
				"ip", clientIP,
			).Infof("HTTP %s %s", method, path)
		}
	}
}
