package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request
		gin.DefaultWriter.Write([]byte(
			gin.LogFormatterParams{
				Request:    c.Request,
				TimeStamp:  time.Now(),
				Latency:    duration,
				ClientIP:   c.ClientIP(),
				Method:     c.Request.Method,
				StatusCode: c.Writer.Status(),
				Path:       c.Request.URL.Path,
			}.String(),
		))
	}
}

