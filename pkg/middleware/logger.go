package middleware

import (
	"time"

	"mastercom-service/pkg/logger"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func Logger(logger *logger.DatadogLogger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP Request", logrus.Fields{
			"client_ip":    param.ClientIP,
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
			"method":       param.Method,
			"path":         param.Path,
			"protocol":     param.Request.Proto,
			"status_code":  param.StatusCode,
			"latency":      param.Latency,
			"user_agent":   param.Request.UserAgent(),
			"error":        param.ErrorMessage,
		})
		return ""
	})
}
