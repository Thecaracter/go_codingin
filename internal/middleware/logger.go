package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logger.WithFields(logrus.Fields{
			"method":     reqMethod,
			"uri":        reqUri,
			"status":     statusCode,
			"latency":    latencyTime,
			"client_ip":  clientIP,
			"user_agent": c.Request.UserAgent(),
		}).Info("HTTP request")
	}
}
