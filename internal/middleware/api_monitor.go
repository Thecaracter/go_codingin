package middleware

import (
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

func APIMonitorMiddleware(apiLogRepo repositories.APILogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Get user ID if authenticated
		var userID *uint
		if uid, exists := c.Get("user_id"); exists {
			id := uid.(uint)
			userID = &id
		}

		// Process request
		c.Next()

		// Calculate response time
		responseTime := time.Since(startTime).Milliseconds()

		// Create API log
		apiLog := &models.APILog{
			UserID:         userID,
			Method:         c.Request.Method,
			Endpoint:       c.Request.URL.Path,
			StatusCode:     c.Writer.Status(),
			ResponseTimeMs: int(responseTime),
			IPAddress:      c.ClientIP(),
			UserAgent:      c.Request.UserAgent(),
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			apiLog.ErrorMessage = c.Errors.String()
		}

		// Save to database (async, don't block response)
		go apiLogRepo.Create(apiLog)
	}
}
